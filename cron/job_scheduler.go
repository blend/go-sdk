package cron

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/blend/go-sdk/stringutil"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/mathutil"
)

// NewJobScheduler returns a job scheduler for a given job.
func NewJobScheduler(job Job, options ...JobSchedulerOption) *JobScheduler {
	js := &JobScheduler{
		Latch: async.NewLatch(),
		Job:   job,
	}

	if typed, ok := job.(JobConfigProvider); ok {
		js.Config = typed.JobConfig()
	}

	if typed, ok := job.(ScheduleProvider); ok {
		js.Schedule = typed.Schedule()
	}

	if typed, ok := job.(DescriptionProvider); ok {
		js.DescriptionProvider = typed.Description
	} else {
		js.DescriptionProvider = func() string { return js.Config.Description }
	}

	if typed, ok := job.(LabelsProvider); ok {
		js.LabelsProvider = typed.Labels
	} else {
		js.LabelsProvider = func() map[string]string { return js.Config.Labels }
	}

	if typed, ok := job.(TimeoutProvider); ok {
		js.TimeoutProvider = typed.Timeout
	} else {
		js.TimeoutProvider = func() time.Duration { return js.Config.TimeoutOrDefault() }
	}

	if typed, ok := job.(ShutdownGracePeriodProvider); ok {
		js.ShutdownGracePeriodProvider = typed.ShutdownGracePeriod
	} else {
		js.ShutdownGracePeriodProvider = func() time.Duration { return js.Config.ShutdownGracePeriodOrDefault() }
	}

	if typed, ok := job.(DisabledProvider); ok {
		js.DisabledProvider = typed.Disabled
	} else {
		js.DisabledProvider = func() bool { return js.Config.DisabledOrDefault() }
	}

	if typed, ok := job.(HistoryDisabledProvider); ok {
		js.HistoryDisabledProvider = typed.HistoryDisabled
	} else {
		js.HistoryDisabledProvider = func() bool { return js.Config.HistoryDisabledOrDefault() }
	}

	if typed, ok := job.(HistoryPersistenceDisabledProvider); ok {
		js.HistoryPersistenceDisabledProvider = typed.HistoryPersistenceDisabled
	} else {
		js.HistoryPersistenceDisabledProvider = func() bool { return js.Config.HistoryPersistenceDisabledOrDefault() }
	}

	if typed, ok := job.(HistoryMaxCountProvider); ok {
		js.HistoryMaxCountProvider = typed.HistoryMaxCount
	} else {
		js.HistoryMaxCountProvider = func() int { return js.Config.HistoryMaxCountOrDefault() }
	}

	if typed, ok := job.(HistoryMaxAgeProvider); ok {
		js.HistoryMaxAgeProvider = typed.HistoryMaxAge
	} else {
		js.HistoryMaxAgeProvider = func() time.Duration { return js.Config.HistoryMaxAgeOrDefault() }
	}

	if typed, ok := job.(ShouldSkipLoggerListenersProvider); ok {
		js.ShouldSkipLoggerListenersProvider = typed.ShouldSkipLoggerListeners
	} else {
		js.ShouldSkipLoggerListenersProvider = func() bool { return js.Config.ShouldSkipLoggerListenersOrDefault() }
	}

	if typed, ok := job.(ShouldSkipLoggerOutputProvider); ok {
		js.ShouldSkipLoggerOutputProvider = typed.ShouldSkipLoggerOutput
	} else {
		js.ShouldSkipLoggerOutputProvider = func() bool { return js.Config.ShouldSkipLoggerOutputOrDefault() }
	}

	if typed, ok := job.(HistoryProvider); ok {
		js.HistoryPersistProvider = typed.PersistHistory
		js.HistoryRestoreProvider = typed.RestoreHistory
	}

	for _, option := range options {
		option(js)
	}

	return js
}

// JobScheduler is a job instance.
type JobScheduler struct {
	sync.Mutex
	Latch *async.Latch

	Job    Job
	Config JobConfig
	Tracer Tracer
	Log    logger.Log

	Schedule    Schedule
	Disabled    bool
	NextRuntime time.Time
	Current     *JobInvocation
	Last        *JobInvocation
	History     []JobInvocation

	DescriptionProvider                func() string
	LabelsProvider                     func() map[string]string
	DisabledProvider                   func() bool
	TimeoutProvider                    func() time.Duration
	ShutdownGracePeriodProvider        func() time.Duration
	ShouldSkipLoggerListenersProvider  func() bool
	ShouldSkipLoggerOutputProvider     func() bool
	HistoryDisabledProvider            func() bool
	HistoryPersistenceDisabledProvider func() bool
	HistoryMaxCountProvider            func() int
	HistoryMaxAgeProvider              func() time.Duration

	HistoryRestoreProvider func(context.Context) ([]JobInvocation, error)
	HistoryPersistProvider func(context.Context, []JobInvocation) error
}

// Name returns the job name.
func (js *JobScheduler) Name() string {
	return js.Job.Name()
}

// Description returns the description.
func (js *JobScheduler) Description() string {
	return js.DescriptionProvider()
}

// Labels returns the job labels, including
// automatically added ones like `name`.
func (js *JobScheduler) Labels() map[string]string {
	output := map[string]string{
		"name":   stringutil.Slugify(js.Name()),
		"state":  string(js.State()),
		"active": fmt.Sprintf("%v", !js.Idle()),
	}
	if js.Last != nil {
		output["last"] = stringutil.Slugify(string(js.Last.State))
	}
	if js.LabelsProvider != nil {
		for key, value := range js.LabelsProvider() {
			output[key] = stringutil.Slugify(value)
		}
	}
	return output
}

// State returns the job scheduler state.
func (js *JobScheduler) State() JobSchedulerState {
	if js.Latch.IsStarted() {
		return JobSchedulerStateRunning
	}
	if js.Latch.IsStopped() {
		return JobSchedulerStateStopped
	}
	return JobSchedulerStateUnknown
}

// Status returns the job scheduler status.
func (js *JobScheduler) Status() JobSchedulerStatus {
	status := JobSchedulerStatus{
		Name:                       js.Name(),
		State:                      js.State(),
		Labels:                     js.Labels(),
		Disabled:                   !js.Enabled(),
		NextRuntime:                js.NextRuntime,
		Timeout:                    js.TimeoutProvider(),
		Current:                    js.Current,
		Last:                       js.Last,
		Stats:                      js.Stats(),
		HistoryDisabled:            js.HistoryDisabledProvider(),
		HistoryPersistenceDisabled: js.HistoryPersistenceDisabledProvider(),
		HistoryMaxCount:            js.HistoryMaxCountProvider(),
		HistoryMaxAge:              js.HistoryMaxAgeProvider(),
	}
	if typed, ok := js.Schedule.(fmt.Stringer); ok {
		status.Schedule = typed.String()
	}
	return status
}

// Stats returns job stats.
func (js *JobScheduler) Stats() JobSchedulerStats {
	output := JobSchedulerStats{
		RunsTotal: len(js.History),
	}
	var elapsedTimes []time.Duration

	for _, ji := range js.History {
		switch ji.State {
		case JobInvocationStateComplete:
			output.RunsSuccessful++
		case JobInvocationStateFailed:
			output.RunsFailed++
		case JobInvocationStateCancelled:
			if !ji.Timeout.IsZero() {
				output.RunsTimedOut++
			} else {
				output.RunsCancelled++
			}
		}

		elapsedTimes = append(elapsedTimes, ji.Elapsed)
		if ji.Elapsed > output.ElapsedMax {
			output.ElapsedMax = ji.Elapsed
		}

		if ji.Output != nil {
			output.OutputBytes += len(ji.Output.Bytes())
		}
	}
	if output.RunsTotal > 0 {
		output.SuccessRate = float64(output.RunsSuccessful) / float64(output.RunsTotal)
	}
	output.Elapsed50th = mathutil.PercentileOfDuration(elapsedTimes, 50.0)
	output.Elapsed95th = mathutil.PercentileOfDuration(elapsedTimes, 95.0)
	return output
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
	stopped := js.Latch.NotifyStopped()

	js.infof("scheduler stopping")
	// signal we are stopping.
	js.Latch.Stopping()

	ctx, cancel := js.createContextWithTimeout(context.Background(), js.ShutdownGracePeriodProvider())
	defer cancel()
	js.cancelJobInvocation(ctx, js.Current)
	js.PersistHistory(ctx)

	<-stopped
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
	js.Disabled = false
	if js.Log != nil && !js.ShouldSkipLoggerListenersProvider() {
		event := NewEvent(FlagEnabled, js.Name())
		js.Log.Trigger(js.logEventContext(context.Background()), event)
	}
	if typed, ok := js.Job.(OnEnabledReceiver); ok {
		typed.OnEnabled(context.Background())
	}
}

// Disable sets the job as disabled.
func (js *JobScheduler) Disable() {
	js.Disabled = true
	if js.Log != nil && !js.ShouldSkipLoggerListenersProvider() {
		event := NewEvent(FlagDisabled, js.Name())
		js.Log.Trigger(js.logEventContext(context.Background()), event)
	}
	if typed, ok := js.Job.(OnDisabledReceiver); ok {
		typed.OnDisabled(context.Background())
	}
}

// Cancel stops all running invocations.
func (js *JobScheduler) Cancel() error {
	if js.Current == nil {
		js.debugf("job cancellation; not running")
		return nil
	}
	gracePeriod := js.ShutdownGracePeriodProvider()
	if gracePeriod > 0 {
		js.debugf("job cancellation; cancelling with %v grace period", gracePeriod)
		ctx, cancel := js.createContextWithTimeout(context.Background(), js.ShutdownGracePeriodProvider())
		defer cancel()

		js.cancelJobInvocation(ctx, js.Current)
	}
	js.debugf("job cancellation; cancelling immediately")
	js.Current.Cancel()
	return nil
}

// RunLoop is the main scheduler loop.
// it alarms on the next runtime and forks a new routine to run the job.
// It can be aborted with the scheduler's async.Latch.
func (js *JobScheduler) RunLoop() {
	js.Latch.Started()
	defer func() {
		js.Latch.Stopped()
	}()

	if js.Schedule != nil {
		js.NextRuntime = js.Schedule.Next(js.NextRuntime)
	}
	// if the schedule returns a zero timestamp
	// it should be interpretted as *not* to automatically
	// schedule the job to be run.
	if js.NextRuntime.IsZero() {
		return
	}

	// this references the underlying js.Latch
	// it returns the current latch signal for stopping *before*
	// the job kicks off.
	notifyStopping := js.Latch.NotifyStopping()

	for {
		if js.NextRuntime.IsZero() {
			return
		}

		runAt := time.After(js.NextRuntime.UTC().Sub(Now()))
		select {
		case <-runAt:
			// if the job is enabled
			// and there isn't another instance running
			if js.CanBeScheduled() {
				// start the job invocation
				go js.Run()
			}

			// set up the next runtime.
			js.NextRuntime = js.Schedule.Next(js.NextRuntime)

		case <-notifyStopping:
			// note: we bail hard here
			// because the job executions in flight are
			// handled by the context cancellation.
			return
		}
	}
}

// RunAsync starts a job invocation with a context.Background() as
// the root context.
func (js *JobScheduler) RunAsync() (*JobInvocation, error) {
	return js.RunAsyncContext(context.Background())
}

// RunAsyncContext starts a job invocation with a given context.
func (js *JobScheduler) RunAsyncContext(ctx context.Context) (*JobInvocation, error) {
	// if there is already another instance running
	if !js.Idle() {
		return nil, ex.New(ErrJobAlreadyRunning, ex.OptMessagef("job: %s", js.Name()))
	}

	timeout := js.TimeoutProvider()

	// create a job invocation, or a record of each
	// individual execution of a job.
	ji := NewJobInvocation(js.Name())
	ji.Context, ji.Cancel = js.createContextWithTimeout(ctx, timeout)

	if timeout > 0 {
		ji.Timeout = ji.Started.Add(timeout)
	}
	js.addCurrent(ji)

	var err error
	var tf TraceFinisher
	// load the job invocation into the context for the job invocation.
	// this will let us pull the job invocation off the context
	// within the job action.
	ji.Context = WithJobInvocation(ji.Context, ji)

	go func() {
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

			js.finishCurrent(ji)
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
		case err = <-js.safeBackgroundExec(ji.Context):
			return
		}
	}()
	return ji, nil
}

// Run forces the job to run.
// This call will block.
func (js *JobScheduler) Run() {
	ji, err := js.RunAsync()
	if err != nil {
		return
	}
	<-ji.Done
}

// RunContext runs a job with a given context as the root context.
func (js *JobScheduler) RunContext(ctx context.Context) {
	ji, err := js.RunAsyncContext(ctx)
	if err != nil {
		return
	}
	<-ji.Done
}

//
// exported utility methods
//

// JobInvocation returns an invocation by id.
func (js *JobScheduler) JobInvocation(id string) *JobInvocation {
	js.Lock()
	defer js.Unlock()

	if js.Current != nil && js.Current.ID == id {
		return js.Current
	}
	if js.Last != nil && js.Last.ID == id {
		return js.Last
	}
	for _, ji := range js.History {
		if ji.ID == id {
			return &ji
		}
	}
	return nil
}

// RestoreHistory calls the persist handler if it's set.
func (js *JobScheduler) RestoreHistory(ctx context.Context) error {
	if js.HistoryPersistenceDisabledProvider() {
		return nil
	}
	if js.HistoryRestoreProvider != nil {
		js.Lock()
		defer js.Unlock()
		var err error
		if js.History, err = js.HistoryRestoreProvider(ctx); err != nil {
			return js.error(err)
		}
		if len(js.History) > 0 {
			js.Last = &js.History[len(js.History)-1]
		}
	}
	return nil
}

// PersistHistory calls the persist handler if it's set.
func (js *JobScheduler) PersistHistory(ctx context.Context) error {
	if js.HistoryPersistenceDisabledProvider() {
		return nil
	}

	if js.HistoryPersistProvider != nil {
		js.Lock()
		defer js.Unlock()

		historyCopy := make([]JobInvocation, len(js.History))
		copy(historyCopy, js.History)
		if err := js.HistoryPersistProvider(ctx, historyCopy); err != nil {
			return js.error(err)
		}
	}
	return nil
}

// CanBeScheduled returns if a job will be triggered automatically
// and isn't already in flight and set to be serial.
func (js *JobScheduler) CanBeScheduled() bool {
	return js.Enabled() && js.Idle()
}

// Enabled returns if the job is explicitly disabled,
// otherwise it checks if the job has an EnabledProvider
// returns the result of that.
// It returns true if there is no provider set.
func (js *JobScheduler) Enabled() bool {
	if js.Disabled {
		return false
	}

	if js.DisabledProvider != nil {
		if js.DisabledProvider() {
			return false
		}
	}

	return true
}

// Idle returns if the job is not currently running.
func (js *JobScheduler) Idle() (isIdle bool) {
	js.Lock()
	isIdle = js.Current == nil
	js.Unlock()
	return
}

//
// utility functions
//

func (js *JobScheduler) finishCurrent(ji *JobInvocation) {
	js.Lock()
	defer js.Unlock()

	if !js.HistoryDisabledProvider() {
		js.History = append(js.cullHistory(), *ji)
	}
	js.Current = nil
	js.Last = ji
	close(ji.Done)
}

func (js *JobScheduler) addCurrent(ji *JobInvocation) {
	js.Lock()
	js.Current = ji
	js.Unlock()
}

func (js *JobScheduler) cullHistory() []JobInvocation {
	count := len(js.History)
	maxCount := js.HistoryMaxCountProvider()
	maxAge := js.HistoryMaxAgeProvider()

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

func (js *JobScheduler) cancelJobInvocation(ctx context.Context, ji *JobInvocation) {
	deadlinePoll := time.Tick(500 * time.Millisecond)
	for {
		if ji == nil || ji.State != JobInvocationStateRunning {
			return
		}
		js.debugf("job cancellation; waiting for cancellation for invocation `%s`", ji.ID)
		select {
		case <-ctx.Done():
			ji.Cancel()
			return
		case <-deadlinePoll:
		}
	}
}

func (js *JobScheduler) safeBackgroundExec(ctx context.Context) chan error {
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

func (js *JobScheduler) createContextWithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout > 0 {
		return context.WithTimeout(ctx, timeout)
	}
	return context.WithCancel(ctx)
}

func (js *JobScheduler) onStart(ctx context.Context, ji *JobInvocation) {
	if js.Log != nil && !js.ShouldSkipLoggerListenersProvider() {
		js.logTrigger(js.logEventContext(ctx), NewEvent(FlagStarted, ji.JobName, OptEventJobInvocation(ji.ID)))
	}
	if typed, ok := js.Job.(OnStartReceiver); ok {
		typed.OnStart(ctx)
	}
}

func (js *JobScheduler) onCancelled(ctx context.Context, ji *JobInvocation) {
	ji.State = JobInvocationStateCancelled

	if js.Log != nil && !js.ShouldSkipLoggerListenersProvider() {
		js.logTrigger(js.logEventContext(ctx), NewEvent(FlagCancelled, ji.JobName, OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed)))
	}
	if typed, ok := js.Job.(OnCancellationReceiver); ok {
		typed.OnCancellation(ctx)
	}
}

func (js *JobScheduler) onComplete(ctx context.Context, ji *JobInvocation) {
	ji.State = JobInvocationStateComplete

	if js.Log != nil && !js.ShouldSkipLoggerListenersProvider() {
		js.logTrigger(js.logEventContext(ctx), NewEvent(FlagComplete, ji.JobName, OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed)))
	}
	if typed, ok := js.Job.(OnCompleteReceiver); ok {
		typed.OnComplete(ctx)
	}

	if js.Last != nil && js.Last.Err != nil {
		js.logTrigger(js.logEventContext(ctx), NewEvent(FlagFixed, ji.JobName, OptEventElapsed(ji.Elapsed)))
		if typed, ok := js.Job.(OnFixedReceiver); ok {
			typed.OnFixed(ctx)
		}
	}
}

func (js *JobScheduler) onFailure(ctx context.Context, ji *JobInvocation) {
	ji.State = JobInvocationStateFailed

	if js.Log != nil && !js.ShouldSkipLoggerListenersProvider() {
		js.logTrigger(js.logEventContext(ctx), NewEvent(FlagFailed, ji.JobName, OptEventErr(ji.Err), OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed)))
	}
	if ji.Err != nil {
		js.error(ji.Err)
	}
	if typed, ok := js.Job.(OnFailureReceiver); ok {
		typed.OnFailure(ctx)
	}
	if js.Last != nil && js.Last.Err == nil {
		if js.Log != nil {
			js.logTrigger(js.logEventContext(ctx), NewEvent(FlagBroken, ji.JobName, OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed)))
		}
		if typed, ok := js.Job.(OnBrokenReceiver); ok {
			typed.OnBroken(ctx)
		}
	}
}

//
// logging helpers
//

func (js *JobScheduler) logEventContext(parent context.Context) context.Context {
	if js.ShouldSkipLoggerOutputProvider() {
		return logger.WithSkipWrite(parent)
	}
	return parent
}

func (js *JobScheduler) logTrigger(ctx context.Context, e logger.Event) {
	if js.Log == nil {
		return
	}
	js.Log.WithPath(js.Name()).Trigger(ctx, e)
}

func (js *JobScheduler) error(err error) error {
	if js.Log == nil {
		return err
	}
	js.Log.WithPath(js.Name()).Error(err)
	return err
}

func (js *JobScheduler) debugf(format string, args ...interface{}) {
	if js.Log == nil {
		return
	}
	js.Log.WithPath(js.Name()).Debugf(format, args...)
}

func (js *JobScheduler) infof(format string, args ...interface{}) {
	if js.Log == nil {
		return
	}
	js.Log.WithPath(js.Name()).Infof(format, args...)
}
