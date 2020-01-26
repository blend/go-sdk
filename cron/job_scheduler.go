package cron

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/ref"
	"github.com/blend/go-sdk/stringutil"
)

// NewJobScheduler returns a job scheduler for a given job.
func NewJobScheduler(job Job, options ...JobSchedulerOption) *JobScheduler {
	js := &JobScheduler{
		Latch: async.NewLatch(),
		Job:   job,
	}
	if typed, ok := job.(ScheduleProvider); ok {
		js.JobSchedule = typed.Schedule()
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

	Job          Job
	JobConfig    JobConfig
	JobSchedule  Schedule
	JobLifecycle JobLifecycle

	Tracer Tracer
	Log    logger.Log

	NextRuntime time.Time
	Current     *JobInvocation
	Last        *JobInvocation
}

// Name returns the job name.
func (js *JobScheduler) Name() string {
	return js.Job.Name()
}

// Config returns the job config provided by a job or an empty config.
func (js *JobScheduler) Config() JobConfig {
	if typed, ok := js.Job.(ConfigProvider); ok {
		return typed.Config()
	}
	return js.JobConfig
}

// Lifecycle returns job lifecycle steps or an empty set.
func (js *JobScheduler) Lifecycle() JobLifecycle {
	if typed, ok := js.Job.(LifecycleProvider); ok {
		return typed.Lifecycle()
	}
	return js.JobLifecycle
}

// Description returns the description.
func (js *JobScheduler) Description() string {
	return js.Config().Description
}

// Disabled returns if the job is disabled or not.
func (js *JobScheduler) Disabled() bool {
	if js.JobConfig.Disabled != nil {
		return *js.JobConfig.Disabled
	}
	return js.Config().DisabledOrDefault()
}

// Labels returns the job labels, including
// automatically added ones like `name`.
func (js *JobScheduler) Labels() map[string]string {
	output := map[string]string{
		"name":      stringutil.Slugify(js.Name()),
		"scheduler": string(js.State()),
		"active":    fmt.Sprint(!js.IsIdle()),
		"enabled":   fmt.Sprint(!js.Disabled()),
	}
	if js.Last != nil {
		output["last"] = stringutil.Slugify(string(js.Last.Status))
	}
	for key, value := range js.Config().Labels {
		output[key] = value
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

// Start starts the scheduler.
// This call blocks.
func (js *JobScheduler) Start() error {
	if !js.Latch.CanStart() {
		return async.ErrCannotStart
	}
	js.Latch.Starting()
	js.RunLoop() // BIG NOTE; this call *should* block
	return nil
}

// Stop stops the scheduler.
func (js *JobScheduler) Stop() error {
	js.Lock()
	defer js.Unlock()

	if !js.Latch.CanStop() {
		return async.ErrCannotStop
	}

	ctx := js.withLogContext(context.Background(), nil)
	js.Latch.Stopping()

	if js.Current != nil {
		gracePeriod := js.Config().ShutdownGracePeriodOrDefault()
		if gracePeriod > 0 {
			var cancel func()
			ctx, cancel = js.withTimeout(ctx, gracePeriod)
			defer cancel()
			js.waitCancelled(ctx, js.Current)
		}
		js.Current.Cancel()
	}

	<-js.Latch.NotifyStopped()
	js.Latch.Reset()
	js.NextRuntime = Zero
	return nil
}

// OnLoad triggers the on load even on the job lifecycle handler.
func (js *JobScheduler) OnLoad(ctx context.Context) error {
	ctx = js.withLogContext(ctx, nil)
	if js.Lifecycle().OnLoad != nil {
		if err := js.Lifecycle().OnLoad(ctx); err != nil {
			return err
		}
	}
	return nil
}

// OnUnload triggers the on unload even on the job lifecycle handler.
func (js *JobScheduler) OnUnload(ctx context.Context) error {
	ctx = js.withLogContext(ctx, nil)
	if js.Lifecycle().OnUnload != nil {
		return js.Lifecycle().OnUnload(ctx)
	}
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
	ctx := js.withLogContext(context.Background(), nil)
	js.JobConfig.Disabled = ref.Bool(false)
	if lifecycle := js.Lifecycle(); lifecycle.OnEnabled != nil {
		lifecycle.OnEnabled(ctx)
	}
	if js.Log != nil && !js.Config().ShouldSkipLoggerListenersOrDefault() {
		js.Log.Trigger(ctx, NewEvent(FlagEnabled, js.Name()))
	}
}

// Disable sets the job as disabled.
func (js *JobScheduler) Disable() {
	ctx := js.withLogContext(context.Background(), nil)
	js.JobConfig.Disabled = ref.Bool(true)
	if lifecycle := js.Lifecycle(); lifecycle.OnDisabled != nil {
		lifecycle.OnDisabled(ctx)
	}
	if js.Log != nil && !js.Config().ShouldSkipLoggerListenersOrDefault() {
		js.Log.Trigger(ctx, NewEvent(FlagDisabled, js.Name()))
	}
}

// Cancel stops all running invocations.
func (js *JobScheduler) Cancel() error {
	if js.Current == nil {
		return nil
	}
	gracePeriod := js.Config().ShutdownGracePeriodOrDefault()
	if gracePeriod > 0 {
		ctx, cancel := js.withTimeout(context.Background(), gracePeriod)
		defer cancel()
		js.waitCancelled(ctx, js.Current)
	}
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
		js.Latch.Reset()
	}()

	ctx := js.withLogContext(context.Background(), nil)
	js.debugf(ctx, "RunLoop: entered running state")

	if js.JobSchedule != nil {
		js.NextRuntime = js.JobSchedule.Next(js.NextRuntime)
	}

	// if the schedule returns a zero timestamp
	// it should be interpretted as *not* to automatically
	// schedule the job to be run.
	// The run loop will return and the job scheduler will be interpretted as stopped.
	if js.NextRuntime.IsZero() {
		js.debugf(ctx, "RunLoop: next runtime is unset, returning")
		return
	}

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
				if _, err := js.RunAsync(); err != nil {
					js.error(ctx, err)
				}
			}

			// set up the next runtime.
			if js.JobSchedule != nil {
				js.NextRuntime = js.JobSchedule.Next(js.NextRuntime)
			} else {
				js.NextRuntime = Zero
			}

		case <-js.Latch.NotifyStopping():
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
	if !js.IsIdle() {
		return nil, ex.New(ErrJobAlreadyRunning, ex.OptMessagef("job: %s", js.Name()))
	}

	timeout := js.Config().TimeoutOrDefault()

	// create a job invocation, or a record of each
	// individual execution of a job.
	ji := NewJobInvocation(js.Name())
	ji.Context = js.withLogContext(ctx, ji)
	ji.Parameters = MergeJobParameterValues(js.Config().ParameterValues, GetJobParameterValues(ji.Context)) // pull the parameters off the calling context.
	ji.Context, ji.Cancel = js.withTimeout(ji.Context, timeout)
	ji.Context = WithJobParameterValues(ji.Context, ji.Parameters)
	ji.Context = WithJobInvocation(ji.Context, ji) // this is confusing but we need to do it
	js.setCurrent(ji)

	var err error
	var tracer TraceFinisher
	go func() {
		// this defer runs all cleanup actions
		// it recovers panics
		// it cancels the timeout (if relevant)
		// it rotates the current and last references
		// it fires lifecycle events
		defer func() {
			if ji.Cancel != nil {
				ji.Cancel()
			}
			if tracer != nil {
				tracer.Finish(ji.Context, err)
			}

			if err != nil && IsJobCancelled(err) {
				js.onJobCancelled(ji.Context, ji)
			} else if err != nil {
				ji.Err = err
				js.onJobError(ji.Context, ji)
			} else {
				js.onJobSuccess(ji.Context, ji)
			}
			js.onJobComplete(ji.Context, ji)
			js.setLast(ji)
		}()

		if js.Tracer != nil {
			ji.Context, tracer = js.Tracer.Start(ji.Context)
		}
		js.onJobBegin(ji.Context, ji)

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

// CanBeScheduled returns if a job will be triggered automatically
// and isn't already in flight and set to be serial.
func (js *JobScheduler) CanBeScheduled() bool {
	return !js.Disabled() && js.IsIdle()
}

// IsIdle returns if the job is not currently running.
func (js *JobScheduler) IsIdle() (isIdle bool) {
	js.Lock()
	isIdle = js.Current == nil
	js.Unlock()
	return
}

//
// utility functions
//

func (js *JobScheduler) setLast(ji *JobInvocation) {
	js.Lock()
	defer js.Unlock()

	js.Current = nil
	js.Last = ji
}

func (js *JobScheduler) setCurrent(ji *JobInvocation) {
	js.Lock()
	js.Current = ji
	js.Unlock()
}

func (js *JobScheduler) waitCancelled(ctx context.Context, ji *JobInvocation) {
	deadlinePoll := time.Tick(500 * time.Millisecond)
	for {
		if ji == nil || ji.Status != JobInvocationStatusRunning {
			return
		}
		select {
		case <-ctx.Done():
			ji.Cancel()
			return
		case <-deadlinePoll:
			continue
		}
	}
}

func (js *JobScheduler) safeBackgroundExec(ctx context.Context) chan error {
	errors := make(chan error, 2)
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

func (js *JobScheduler) withTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout > 0 {
		return context.WithTimeout(ctx, timeout)
	}
	return context.WithCancel(ctx)
}

// job lifecycle hooks

func (js *JobScheduler) onJobBegin(ctx context.Context, ji *JobInvocation) {
	defer func() {
		if r := recover(); r != nil {
			js.error(ctx, ex.New(r))
		}
	}()

	ji.Started = time.Now().UTC()
	ji.Status = JobInvocationStatusRunning

	if lifecycle := js.Lifecycle(); lifecycle.OnBegin != nil {
		lifecycle.OnBegin(ctx)
	}
	if js.Log != nil && !js.Config().ShouldSkipLoggerListenersOrDefault() {
		js.logTrigger(ctx, NewEvent(FlagBegin, ji.JobName, OptEventJobInvocation(ji.ID)))
	}
}

func (js *JobScheduler) onJobComplete(ctx context.Context, ji *JobInvocation) {
	defer func() {
		if r := recover(); r != nil {
			js.error(ctx, ex.New(r))
		}
		close(ji.Done)
	}()

	// this must be set here so the handlers pick it up
	ji.Complete = time.Now().UTC()
	if lifecycle := js.Lifecycle(); lifecycle.OnComplete != nil {
		lifecycle.OnComplete(ctx)
	}
	if js.Log != nil && !js.Config().ShouldSkipLoggerListenersOrDefault() {
		js.logTrigger(ctx, NewEvent(FlagComplete, ji.JobName, OptEventElapsed(ji.Elapsed())))
	}
}

func (js *JobScheduler) onJobCancelled(ctx context.Context, ji *JobInvocation) {
	defer func() {
		if r := recover(); r != nil {
			js.error(ctx, ex.New(r))
		}
	}()

	ji.Status = JobInvocationStatusCancelled
	if lifecycle := js.Lifecycle(); lifecycle.OnCancellation != nil {
		lifecycle.OnCancellation(ctx)
	}
	if js.Log != nil && !js.Config().ShouldSkipLoggerListenersOrDefault() {
		js.logTrigger(ctx, NewEvent(FlagCancelled, ji.JobName, OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed())))
	}
}

func (js *JobScheduler) onJobSuccess(ctx context.Context, ji *JobInvocation) {
	defer func() {
		if r := recover(); r != nil {
			js.error(ctx, ex.New(r))
		}
	}()

	ji.Status = JobInvocationStatusSuccess
	if lifecycle := js.Lifecycle(); lifecycle.OnSuccess != nil {
		lifecycle.OnSuccess(ctx)
	}
	if js.Log != nil && !js.Config().ShouldSkipLoggerListenersOrDefault() {
		js.logTrigger(ctx, NewEvent(FlagSuccess, ji.JobName, OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed())))
	}
	if js.Last != nil && js.Last.Err != nil {
		if lifecycle := js.Lifecycle(); lifecycle.OnFixed != nil {
			lifecycle.OnFixed(ctx)
		}
		if js.Log != nil && !js.Config().ShouldSkipLoggerListenersOrDefault() {
			js.logTrigger(ctx, NewEvent(FlagFixed, ji.JobName, OptEventElapsed(ji.Elapsed())))
		}
	}
}

func (js *JobScheduler) onJobError(ctx context.Context, ji *JobInvocation) {
	defer func() {
		if r := recover(); r != nil {
			js.error(ctx, ex.New(r))
		}
	}()

	ji.Status = JobInvocationStatusErrored

	if lifecycle := js.Lifecycle(); lifecycle.OnError != nil {
		lifecycle.OnError(ctx)
	}
	if js.Log != nil && !js.Config().ShouldSkipLoggerListenersOrDefault() {
		js.logTrigger(ctx, NewEvent(FlagErrored, ji.JobName, OptEventErr(ji.Err), OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed())))
	}
	if ji.Err != nil {
		js.error(ctx, ji.Err)
	}

	if js.Last != nil && js.Last.Err == nil {
		if lifecycle := js.Lifecycle(); lifecycle.OnBroken != nil {
			lifecycle.OnBroken(ctx)
		}
		if js.Log != nil && !js.Config().ShouldSkipLoggerListenersOrDefault() {
			js.logTrigger(ctx, NewEvent(FlagBroken, ji.JobName, OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed())))
		}
	}
}

//
// logging helpers
//

func (js *JobScheduler) withLogContext(parent context.Context, ji *JobInvocation) context.Context {
	if ji != nil {
		parent = logger.WithPath(parent, js.Name(), ji.ID)
	} else {
		parent = logger.WithPath(parent, js.Name())
	}
	if js.Config().ShouldSkipLoggerOutputOrDefault() {
		parent = logger.WithSkipWrite(parent)
	}
	return parent
}

func (js *JobScheduler) logTrigger(ctx context.Context, e logger.Event) {
	if js.Log == nil {
		return
	}
	js.Log.Trigger(ctx, e)
}

func (js *JobScheduler) debugf(ctx context.Context, format string, args ...interface{}) {
	if js.Log == nil {
		return
	}
	js.Log.WithContext(ctx).Debugf(format, args...)
}

func (js *JobScheduler) error(ctx context.Context, err error) error {
	if js.Log == nil {
		return err
	}
	js.Log.WithContext(ctx).Error(err)
	return err
}
