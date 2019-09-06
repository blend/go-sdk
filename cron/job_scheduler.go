package cron

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
)

// NewJobScheduler returns a job scheduler for a given job.
func NewJobScheduler(job Job, options ...JobSchedulerOption) *JobScheduler {
	js := &JobScheduler{
		Latch: async.NewLatch(),
		Name:  job.Name(),
		Job:   job,
	}

	if typed, ok := job.(DescriptionProvider); ok {
		js.Description = typed.Description()
	}

	if typed, ok := job.(ScheduleProvider); ok {
		js.Schedule = typed.Schedule()
	}

	if typed, ok := job.(TimeoutProvider); ok {
		js.TimeoutProvider = typed.Timeout
	} else {
		js.TimeoutProvider = func() time.Duration { return 0 }
	}

	if typed, ok := job.(ShutdownGracePeriodProvider); ok {
		js.ShutdownGracePeriodProvider = typed.ShutdownGracePeriod
	} else {
		js.ShutdownGracePeriodProvider = func() time.Duration { return 0 }
	}

	if typed, ok := job.(EnabledProvider); ok {
		js.EnabledProvider = typed.Enabled
	} else {
		js.EnabledProvider = func() bool { return DefaultEnabled }
	}

	if typed, ok := job.(SerialProvider); ok {
		js.SerialProvider = typed.Serial
	} else {
		js.SerialProvider = func() bool { return DefaultSerial }
	}

	if typed, ok := job.(ShouldTriggerListenersProvider); ok {
		js.ShouldTriggerListenersProvider = typed.ShouldTriggerListeners
	} else {
		js.ShouldTriggerListenersProvider = func() bool { return DefaultShouldTriggerListeners }
	}

	if typed, ok := job.(ShouldWriteOutputProvider); ok {
		js.ShouldWriteOutputProvider = typed.ShouldWriteOutput
	} else {
		js.ShouldWriteOutputProvider = func() bool { return DefaultShouldWriteOutput }
	}

	if typed, ok := job.(HistoryPersister); ok {
		js.HistoryPersist = typed.HistoryPersist
	}

	for _, option := range options {
		option(js)
	}

	return js
}

// JobScheduler is a job instance.
type JobScheduler struct {
	sync.Mutex   `json:"-"`
	*async.Latch `json:"-"`

	Name        string `json:"name"`
	Description string `json:"description"`
	Job         Job    `json:"-"`

	Config Config     `json:"-"`
	Tracer Tracer     `json:"-"`
	Log    logger.Log `json:"-"`

	// Meta Fields
	Disabled    bool            `json:"disabled"`
	NextRuntime time.Time       `json:"nextRuntime"`
	Current     *JobInvocation  `json:"current"`
	Last        *JobInvocation  `json:"last"`
	History     []JobInvocation `json:"history"`

	Schedule                       Schedule                                     `json:"-"`
	EnabledProvider                func() bool                                  `json:"-"`
	SerialProvider                 func() bool                                  `json:"-"`
	TimeoutProvider                func() time.Duration                         `json:"-"`
	ShutdownGracePeriodProvider    func() time.Duration                         `json:"-"`
	ShouldTriggerListenersProvider func() bool                                  `json:"-"`
	ShouldWriteOutputProvider      func() bool                                  `json:"-"`
	HistoryPersist                 func(context.Context, []JobInvocation) error `json:"-"`
}

// Start starts the scheduler.
// This call blocks.
func (js *JobScheduler) Start() error {
	if !js.Latch.CanStart() {
		return fmt.Errorf("already started")
	}
	js.infof("scheduler starting")
	js.Latch.Starting()
	js.infof("scheduler started")
	js.RunLoop()
	js.infof("scheduler exiting")
	return nil
}

// Stop stops the scheduler.
func (js *JobScheduler) Stop() error {
	if !js.Latch.CanStop() {
		return fmt.Errorf("already stopped")
	}
	// Signal we are stopping.
	js.infof("scheduler stopping")
	js.Latch.Stopping()

	ctx, cancel := js.createContextWithTimeout(js.ShutdownGracePeriodProvider())
	defer cancel()
	js.cancel(ctx)
	js.PersistHistory(ctx)

	<-js.Latch.NotifyStopped()
	js.infof("scheduler stopped")
	return nil
}

// NotifyStarted notifies the job scheduler has started.
func (js *JobScheduler) NotifyStarted() <-chan struct{} {
	return js.Latch.NotifyStarted()
}

// NotifyStopped notifies the job scheduler has stopped.
func (js *JobScheduler) NotifyStopped() <-chan struct{} {
	return js.Latch.NotifyStopped()
}

// Enable sets the job as enabled.
func (js *JobScheduler) Enable() {
	js.Lock()
	defer js.Unlock()

	js.Disabled = false
	if js.Log != nil && js.ShouldTriggerListenersProvider() {
		event := NewEvent(FlagEnabled, js.Name, OptEventWritable(js.ShouldWriteOutputProvider()))
		js.Log.Trigger(context.Background(), event)
	}
	if typed, ok := js.Job.(OnEnabledReceiver); ok {
		typed.OnEnabled(context.Background())
	}
}

// Disable sets the job as disabled.
func (js *JobScheduler) Disable() {
	js.Lock()
	defer js.Unlock()

	js.Disabled = true
	if js.Log != nil && js.ShouldTriggerListenersProvider() {
		event := NewEvent(FlagDisabled, js.Name, OptEventWritable(js.ShouldWriteOutputProvider()))
		js.Log.Trigger(context.Background(), event)
	}
	if typed, ok := js.Job.(OnDisabledReceiver); ok {
		typed.OnDisabled(context.Background())
	}
}

// Cancel stops an execution in process.
func (js *JobScheduler) Cancel() error {
	if js.Current == nil {
		return nil
	}
	gracePeriod := js.ShutdownGracePeriodProvider()
	if gracePeriod > 0 {
		ctx, cancel := js.createContextWithTimeout(js.ShutdownGracePeriodProvider())
		defer cancel()
		return js.cancel(ctx)
	}
	js.Current.Cancel()
	return nil
}

// Cancel stops an execution in process.
func (js *JobScheduler) cancel(ctx context.Context) error {
	deadlinePoll := time.Tick(500 * time.Millisecond)
	for {
		if js.Current == nil {
			return nil
		}
		js.debugf("job in flight cancelling")
		select {
		case <-ctx.Done():
			if js.Current != nil {
				js.debugf("job in flight cancelled")
				js.Current.Cancel()
			}
			return nil
		case <-deadlinePoll:
		}
	}
}

// RunLoop is the main scheduler loop.
// it alarms on the next runtime and forks a new routine to run the job.
// It can be aborted with the scheduler's async.Latch.
func (js *JobScheduler) RunLoop() {
	js.Started()
	defer func() {
		js.Stopped()
	}()

	if js.Schedule != nil {
		js.NextRuntime = js.Schedule.Next(js.NextRuntime)
	}
	if js.NextRuntime.IsZero() {
		return
	}

	var notifyStopping <-chan struct{}
	for {
		if js.NextRuntime.IsZero() {
			return
		}
		// this references the underlying js.Latch
		// it returns the current latch signal for stopping *before*
		// the job kicks off.
		notifyStopping = js.NotifyStopping()
		runAt := time.After(js.NextRuntime.UTC().Sub(Now()))
		select {
		case <-runAt:
			if js.enabled() {
				// start the job
				go js.Run()
			}
			// set up the next runtime.
			js.NextRuntime = js.Schedule.Next(js.NextRuntime)
		case <-notifyStopping:
			// note: we bail hard here
			// because the job executions in flight are
			// responsible for themselves.
			return
		}
	}
}

// Run forces the job to run.
// It checks if the job should be allowed to execute.
// It blocks on the job execution to enforce or clear timeouts.
func (js *JobScheduler) Run() {
	// check if the job can run
	if !js.enabled() {
		return
	}

	timeout := js.TimeoutProvider()

	// create a job invocation, or a record of each
	// individual execution of a job.
	ji := NewJobInvocation(js.Name)
	ji.Context, ji.Cancel = js.createContextWithTimeout(timeout)

	if timeout > 0 {
		ji.Timeout = ji.Started.Add(timeout)
	}
	js.setCurrent(ji)

	var err error
	var tf TraceFinisher

	// load the job invocation into the context
	ji.Context = WithJobInvocation(ji.Context, ji)

	// this defer runs all cleanup actions
	// it recovers panics
	// it cancels the timeout (if relevant)
	// it rotates the current and last references
	// it fires lifecycle events
	defer func() {
		if r := recover(); r != nil {
			err = ex.New(err)
		}
		if ji.Cancel != nil {
			ji.Cancel()
		}
		if tf != nil {
			tf.Finish(ji.Context)
		}

		ji.Finished = Now()
		ji.Elapsed = ji.Finished.Sub(ji.Started)
		ji.Err = err

		if err != nil && IsJobCancelled(err) {
			ji.Cancelled = ji.Finished
			js.onCancelled(ji.Context, ji)
		} else if ji.Err != nil {
			js.onFailure(ji.Context, ji)
		} else {
			js.onComplete(ji.Context, ji)
		}

		js.addHistory(*ji)
		js.setCurrent(nil)
		js.setLast(ji)
		js.PersistHistory(ji.Context)
	}()

	// if the tracer is set, create a trace context
	if js.Tracer != nil {
		ji.Context, tf = js.Tracer.Start(ji.Context)
	}
	// fire the on start event
	js.onStart(ji.Context, ji)

	// check if the job has been canceled
	// or if it's finished.
	select {
	case <-ji.Context.Done():
		err = ErrJobCancelled
		return
	case err = <-js.safeAsyncExec(ji.Context):
		return
	}
}

//
// exported utility methods
//

// GetInvocationByID returns an invocation by id.
func (js *JobScheduler) GetInvocationByID(id string) *JobInvocation {
	for _, ji := range js.History {
		if ji.ID == id {
			return &ji
		}
	}
	return nil
}

// PersistHistory calls the persist handler if it's set.
func (js *JobScheduler) PersistHistory(ctx context.Context) error {
	if js.HistoryPersist != nil {
		if err := js.HistoryPersist(ctx, js.History); err != nil {
			return js.error(err)
		}
	}
	return nil
}

//
// utility functions
//

func (js *JobScheduler) setCurrent(ji *JobInvocation) {
	js.Current = ji
}

func (js *JobScheduler) setLast(ji *JobInvocation) {
	js.Last = ji
}

// safeAsyncExec runs a given job's body and recovers panics.
func (js *JobScheduler) safeAsyncExec(ctx context.Context) chan error {
	errors := make(chan error)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errors <- ex.New(r)
			}
		}()
		errors <- js.Job.Execute(ctx)
	}()
	return errors
}

func (js *JobScheduler) createContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout > 0 {
		return context.WithTimeout(context.Background(), timeout)
	}
	return context.WithCancel(context.Background())
}

// enabled returns if a job can execute.
func (js *JobScheduler) enabled() bool {
	if js.Disabled {
		return false
	}

	if js.EnabledProvider != nil {
		if !js.EnabledProvider() {
			return false
		}
	}

	if js.SerialProvider != nil && js.SerialProvider() {
		if js.Current != nil {
			return false
		}
	}
	return true
}

func (js *JobScheduler) onStart(ctx context.Context, ji *JobInvocation) {
	if js.Log != nil && js.ShouldTriggerListenersProvider() {
		js.trigger(ctx, NewEvent(FlagStarted, ji.JobName, OptEventJobInvocation(ji.ID), OptEventWritable(js.ShouldWriteOutputProvider())))
	}
	if typed, ok := js.Job.(OnStartReceiver); ok {
		typed.OnStart(ctx)
	}
}

func (js *JobScheduler) onCancelled(ctx context.Context, ji *JobInvocation) {
	ji.Status = JobStatusCancelled

	if js.Log != nil && js.ShouldTriggerListenersProvider() {
		js.trigger(ctx, NewEvent(FlagCancelled, ji.JobName, OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed), OptEventWritable(js.ShouldWriteOutputProvider())))
	}
	if typed, ok := js.Job.(OnCancellationReceiver); ok {
		typed.OnCancellation(ctx)
	}
}

func (js *JobScheduler) onComplete(ctx context.Context, ji *JobInvocation) {
	ji.Status = JobStatusComplete

	if js.Log != nil && js.ShouldTriggerListenersProvider() {
		js.trigger(ctx, NewEvent(FlagComplete, ji.JobName, OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed), OptEventWritable(js.ShouldWriteOutputProvider())))
	}
	if typed, ok := js.Job.(OnCompleteReceiver); ok {
		typed.OnComplete(ctx)
	}

	if js.Last != nil && js.Last.Err != nil {
		js.trigger(ctx, NewEvent(FlagFixed, ji.JobName, OptEventElapsed(ji.Elapsed), OptEventWritable(js.ShouldWriteOutputProvider())))
		if typed, ok := js.Job.(OnFixedReceiver); ok {
			typed.OnFixed(ctx)
		}
	}
}

func (js *JobScheduler) onFailure(ctx context.Context, ji *JobInvocation) {
	ji.Status = JobStatusFailed

	if js.Log != nil && js.ShouldTriggerListenersProvider() {
		js.trigger(ctx, NewEvent(FlagFailed, ji.JobName, OptEventErr(ji.Err), OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed), OptEventWritable(js.ShouldWriteOutputProvider())))
	}
	if ji.Err != nil {
		js.error(ji.Err)
	}
	if typed, ok := js.Job.(OnFailureReceiver); ok {
		typed.OnFailure(ctx)
	}
	if js.Last != nil && js.Last.Err == nil {
		if js.Log != nil {
			js.trigger(ctx, NewEvent(FlagBroken, ji.JobName, OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed), OptEventWritable(js.ShouldWriteOutputProvider())))
		}
		if typed, ok := js.Job.(OnBrokenReceiver); ok {
			typed.OnBroken(ctx)
		}
	}
}

func (js *JobScheduler) addHistory(ji JobInvocation) {
	js.History = append(js.cullHistory(), ji)
}

func (js *JobScheduler) cullHistory() []JobInvocation {
	count := len(js.History)
	maxCount := js.Config.HistoryMaxCountOrDefault()
	maxAge := js.Config.HistoryMaxAgeOrDefault()
	now := time.Now().UTC()
	var filtered []JobInvocation
	for index, h := range js.History {
		if maxCount > 0 {
			if index < (count - maxCount) {
				continue
			}
		}
		if maxAge > 0 {
			if now.Sub(h.Started) > maxAge {
				continue
			}
		}
		filtered = append(filtered, h)
	}
	return filtered
}

func (js *JobScheduler) trigger(ctx context.Context, e logger.Event) {
	if js.Log == nil {
		return
	}
	js.Log.WithPath(js.Name).Trigger(ctx, e)
}

func (js *JobScheduler) error(err error) error {
	if js.Log == nil {
		return err
	}
	js.Log.WithPath(js.Name).Error(err)
	return err
}

func (js *JobScheduler) debugf(format string, args ...interface{}) {
	if js.Log == nil {
		return
	}
	js.Log.WithPath(js.Name).Debugf(format, args...)
}

func (js *JobScheduler) infof(format string, args ...interface{}) {
	if js.Log == nil {
		return
	}
	js.Log.WithPath(js.Name).Infof(format, args...)
}
