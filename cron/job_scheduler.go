package cron

import (
	"context"
	"sync"
	"time"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
)

// NewJobScheduler returns a job scheduler for a given job.
func NewJobScheduler(cfg *Config, job Job) *JobScheduler {
	js := &JobScheduler{
		Latch:  &async.Latch{},
		Name:   job.Name(),
		Job:    job,
		Config: cfg,
	}

	if typed, ok := job.(ScheduleProvider); ok {
		js.Schedule = typed.Schedule()
	}

	if typed, ok := job.(TimeoutProvider); ok {
		js.TimeoutProvider = typed.Timeout
	} else {
		js.TimeoutProvider = func() time.Duration { return 0 }
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

	return js
}

// JobScheduler is a job instance.
type JobScheduler struct {
	sync.Mutex `json:"-"`
	Latch      *async.Latch `json:"-"`

	Name string `json:"name"`
	Job  Job    `json:"-"`

	Tracer Tracer              `json:"-"`
	Log    logger.FullReceiver `json:"-"`
	Config *Config             `json:"-"`

	// Meta Fields
	Disabled    bool            `json:"disabled"`
	NextRuntime time.Time       `json:"nextRuntime"`
	Current     *JobInvocation  `json:"current"`
	Last        *JobInvocation  `json:"last"`
	History     []JobInvocation `json:"history"`

	Schedule                       Schedule             `json:"-"`
	EnabledProvider                func() bool          `json:"-"`
	SerialProvider                 func() bool          `json:"-"`
	TimeoutProvider                func() time.Duration `json:"-"`
	ShouldTriggerListenersProvider func() bool          `json:"-"`
	ShouldWriteOutputProvider      func() bool          `json:"-"`
}

// Invocation returns an invocation by id.
func (js *JobScheduler) Invocation(id string) *JobInvocation {
	for _, ji := range js.History {
		if ji.ID == id {
			return &ji
		}
	}
	return nil
}

// WithTracer sets the scheduler tracer.
func (js *JobScheduler) WithTracer(tracer Tracer) *JobScheduler {
	js.Tracer = tracer
	return js
}

// WithLogger sets the scheduler logger.
func (js *JobScheduler) WithLogger(log logger.FullReceiver) *JobScheduler {
	js.Log = log
	return js
}

// Start starts the scheduler.
func (js *JobScheduler) Start() {
	if !js.Latch.CanStart() {
		return
	}
	js.Latch.Starting()
	go js.RunLoop()
	<-js.Latch.NotifyStarted()
}

// Stop stops the scheduler.
func (js *JobScheduler) Stop() {
	if !js.Latch.CanStop() {
		return
	}
	js.Latch.Stopping()
	<-js.Latch.NotifyStopped()
}

// Enable sets the job as enabled.
func (js *JobScheduler) Enable() {
	js.Lock()
	defer js.Unlock()

	js.Disabled = false
	if js.Log != nil && js.ShouldTriggerListenersProvider() {
		event := NewEvent(FlagEnabled, js.Name).
			WithIsWritable(js.ShouldWriteOutputProvider())
		js.Log.Trigger(event)
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
		event := NewEvent(FlagDisabled, js.Name).
			WithIsWritable(js.ShouldWriteOutputProvider())
		js.Log.Trigger(event)
	}
	if typed, ok := js.Job.(OnDisabledReceiver); ok {
		typed.OnDisabled(context.Background())
	}
}

// Cancel stops an execution in process.
func (js *JobScheduler) Cancel() {
	if js.Current != nil {
		js.Current.Cancel()
	}
}

// RunLoop is the main scheduler loop.
// it alarms on the next runtime and forks a new routine to run the job.
// It can be aborted with the scheduler's async.Latch.
func (js *JobScheduler) RunLoop() {
	js.Latch.Started()

	if js.Schedule != nil {
		// sniff the schedule, see if a next runtime is called for (or if the job is on demand).
		js.NextRuntime = js.Schedule.Next(js.NextRuntime)
	}
	if js.NextRuntime.IsZero() {
		js.Latch.Stopped()
		return
	}

	for {
		if js.NextRuntime.IsZero() {
			js.Latch.Stopped()
			return
		}
		runAt := time.After(js.NextRuntime.UTC().Sub(Now()))
		select {
		case <-runAt:
			// start the job
			go js.Run()
			// set up the next runtime.
			js.NextRuntime = js.Schedule.Next(js.NextRuntime)
		case <-js.Latch.NotifyStopping():
			js.Latch.Stopped()
			return
		}
	}
}

// Run forces the job to run.
// It checks if the job should be allowed to execute.
// It blocks on the job execution to enforce or clear timeouts.
func (js *JobScheduler) Run() {
	// check if the job can run
	if !js.canRun() {
		return
	}

	// mark the start time
	start := Now()

	timeout := js.TimeoutProvider()

	// create the root context.
	ctx, cancel := js.createContextWithTimeout(timeout)

	// create a job invocation, or a record of each
	// individual execution of a job.
	ji := JobInvocation{
		ID:      NewJobInvocationID(),
		JobName: js.Name,
		Status:  JobStatusRunning,
		Started: start,
		Context: ctx,
		Cancel:  cancel,
	}
	if timeout > 0 {
		ji.Timeout = start.Add(timeout)
	}
	js.setCurrent(&ji)

	var err error
	var tf TraceFinisher
	// load the job invocation into the context
	ctx = WithJobInvocation(ctx, &ji)

	// this defer runs all cleanup actions
	// it recovers panics
	// it cancels the timeout (if relevant)
	// it rotates the current and last references
	// it fires lifecycle events
	defer func() {
		if r := recover(); r != nil {
			err = exception.New(err)
		}
		cancel()
		if tf != nil {
			tf.Finish(ctx)
		}

		ji.Finished = Now()
		ji.Elapsed = ji.Finished.Sub(ji.Started)
		ji.Err = err

		if err != nil && IsJobCancelled(err) {
			ji.Cancelled = ji.Finished
			js.onCancelled(ctx, &ji)
		} else if ji.Err != nil {
			js.onFailure(ctx, &ji)
		} else {
			js.onComplete(ctx, &ji)
		}

		js.addHistory(ji)
		js.setCurrent(nil)
		js.setLast(&ji)
	}()

	// if the tracer is set, create a trace context
	if js.Tracer != nil {
		ctx, tf = js.Tracer.Start(ctx)
	}
	// fire the on start event
	js.onStart(ctx, &ji)

	// check if the job has been canceled
	// or if it's finished.
	select {
	case <-ctx.Done():
		err = ErrJobCancelled
	case err = <-js.safeAsyncExec(ctx):
	}
}

//
// utility functions
//

func (js *JobScheduler) setCurrent(ji *JobInvocation) {
	js.Lock()
	js.Current = ji
	js.Unlock()
}

func (js *JobScheduler) setLast(ji *JobInvocation) {
	js.Lock()
	js.Last = ji
	js.Unlock()
}

// safeAsyncExec runs a given job's body and recovers panics.
func (js *JobScheduler) safeAsyncExec(ctx context.Context) chan error {
	errors := make(chan error)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errors <- exception.New(r)
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

// canRun returns if a job can execute.
func (js *JobScheduler) canRun() bool {
	js.Lock()
	defer js.Unlock()

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
		event := NewEvent(FlagStarted, ji.JobName).
			WithJobInvocation(ji.ID).
			WithIsWritable(js.ShouldWriteOutputProvider())
		js.Log.Trigger(event)
	}
	if typed, ok := js.Job.(OnStartReceiver); ok {
		typed.OnStart(ctx)
	}
}

func (js *JobScheduler) onCancelled(ctx context.Context, ji *JobInvocation) {
	ji.Status = JobStatusCancelled

	if js.Log != nil && js.ShouldTriggerListenersProvider() {
		event := NewEvent(FlagCancelled, ji.JobName).
			WithJobInvocation(ji.ID).
			WithIsWritable(js.ShouldWriteOutputProvider()).
			WithElapsed(ji.Elapsed)
		js.Log.Trigger(event)
	}
	if typed, ok := js.Job.(OnCancellationReceiver); ok {
		typed.OnCancellation(ctx)
	}
}

func (js *JobScheduler) onComplete(ctx context.Context, ji *JobInvocation) {
	ji.Status = JobStatusComplete

	if js.Log != nil && js.ShouldTriggerListenersProvider() {
		event := NewEvent(FlagComplete, ji.JobName).
			WithJobInvocation(ji.ID).
			WithIsWritable(js.ShouldWriteOutputProvider()).
			WithElapsed(ji.Elapsed)
		js.Log.Trigger(event)
	}
	if typed, ok := js.Job.(OnCompleteReceiver); ok {
		typed.OnComplete(ctx)
	}

	if js.Last != nil && js.Last.Err != nil {
		if js.Log != nil {
			event := NewEvent(FlagFixed, ji.JobName).
				WithIsWritable(js.ShouldWriteOutputProvider()).
				WithElapsed(ji.Elapsed)

			js.Log.Trigger(event)
		}

		if typed, ok := js.Job.(OnFixedReceiver); ok {
			typed.OnFixed(ctx)
		}
	}
}

func (js *JobScheduler) onFailure(ctx context.Context, ji *JobInvocation) {
	ji.Status = JobStatusFailed

	if js.Log != nil && js.ShouldTriggerListenersProvider() {
		event := NewEvent(FlagFailed, ji.JobName).
			WithJobInvocation(ji.ID).
			WithIsWritable(js.ShouldWriteOutputProvider()).
			WithElapsed(ji.Elapsed).
			WithErr(ji.Err)

		js.Log.Trigger(event)
	}
	if ji.Err != nil {
		logger.MaybeError(js.Log, ji.Err)
	}
	if typed, ok := js.Job.(OnFailureReceiver); ok {
		typed.OnFailure(ctx)
	}
	if js.Last != nil && js.Last.Err == nil {
		if js.Log != nil {
			event := NewEvent(FlagBroken, ji.JobName).
				WithJobInvocation(ji.ID).
				WithIsWritable(js.ShouldWriteOutputProvider()).
				WithElapsed(ji.Elapsed)
			js.Log.Trigger(event)
		}

		if typed, ok := js.Job.(OnBrokenReceiver); ok {
			typed.OnBroken(ctx)
		}
	}
}

func (js *JobScheduler) addHistory(ji JobInvocation) {
	js.Lock()
	defer js.Unlock()
	js.History = append(js.cullHistory(), ji)
}

func (js *JobScheduler) cullHistory() []JobInvocation {
	if js.Config == nil {
		return js.History
	}
	count := len(js.History)
	maxCount := js.Config.History.MaxCountOrDefault()
	maxAge := js.Config.History.MaxAgeOrDefault()
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
