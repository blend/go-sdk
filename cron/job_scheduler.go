/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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
		Latch:       async.NewLatch(),
		BaseContext: context.Background(),
		Job:         job,
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
	Latch *async.Latch

	Job          Job
	JobConfig    JobConfig
	JobSchedule  Schedule
	JobLifecycle JobLifecycle

	BaseContext context.Context

	Tracer Tracer
	Log    logger.Log

	NextRuntime time.Time

	currentLock sync.Mutex
	current     *JobInvocation
	lastLock    sync.Mutex
	last        *JobInvocation
}

// Name returns the job name.
func (js *JobScheduler) Name() string {
	return js.Job.Name()
}

// Background returns the job scheduler base context.
//
// It should be used as the root context for _any_ operations.
func (js *JobScheduler) Background() context.Context {
	if js.BaseContext != nil {
		return js.BaseContext
	}
	return context.Background()
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
	if js.Last() != nil {
		output["last"] = stringutil.Slugify(string(js.Last().Status))
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
	js.RunLoop()
	return nil
}

// Stop stops the scheduler.
func (js *JobScheduler) Stop() error {
	if !js.Latch.CanStop() {
		return async.ErrCannotStop
	}

	ctx := js.withBaseContext(js.Background())
	js.Latch.Stopping()

	if current := js.Current(); current != nil {
		gracePeriod := js.Config().ShutdownGracePeriodOrDefault()
		if gracePeriod > 0 {
			var cancel func()
			ctx, cancel = js.withTimeoutOrCancel(ctx, gracePeriod)
			defer cancel()
			js.waitCurrentComplete(ctx)
		}
	}
	if current := js.Current(); current != nil && current.Status == JobInvocationStatusRunning {
		current.Cancel()
	}

	<-js.Latch.NotifyStopped()
	js.Latch.Reset()
	js.NextRuntime = Zero
	return nil
}

// OnLoad triggers the on load even on the job lifecycle handler.
func (js *JobScheduler) OnLoad(ctx context.Context) error {
	ctx = js.withBaseContext(ctx)
	if js.Lifecycle().OnLoad != nil {
		if err := js.Lifecycle().OnLoad(ctx); err != nil {
			return err
		}
	}
	return nil
}

// OnUnload triggers the on unload even on the job lifecycle handler.
func (js *JobScheduler) OnUnload(ctx context.Context) error {
	ctx = js.withBaseContext(ctx)
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
	ctx := js.withBaseContext(js.Background())
	js.JobConfig.Disabled = ref.Bool(false)
	if lifecycle := js.Lifecycle(); lifecycle.OnEnabled != nil {
		lifecycle.OnEnabled(ctx)
	}
	if js.Log != nil && !js.Config().SkipLoggerTrigger {
		js.Log.TriggerContext(ctx, NewEvent(FlagEnabled, js.Name()))
	}
}

// Disable sets the job as disabled.
func (js *JobScheduler) Disable() {
	ctx := js.withBaseContext(js.Background())
	js.JobConfig.Disabled = ref.Bool(true)
	if lifecycle := js.Lifecycle(); lifecycle.OnDisabled != nil {
		lifecycle.OnDisabled(ctx)
	}
	if js.Log != nil && !js.Config().SkipLoggerTrigger {
		js.Log.TriggerContext(ctx, NewEvent(FlagDisabled, js.Name()))
	}
}

// Cancel stops all running invocations.
func (js *JobScheduler) Cancel() error {
	ctx := js.withBaseContext(js.Background())

	if js.Current() == nil {
		logger.MaybeDebugfContext(ctx, js.Log, "cannot cancel; job is not runnning")
		return nil
	}
	gracePeriod := js.Config().ShutdownGracePeriodOrDefault()
	if gracePeriod > 0 {
		ctx, cancel := js.withTimeoutOrCancel(ctx, gracePeriod)
		defer cancel()
		js.waitCurrentComplete(ctx)
	}
	if current := js.Current(); current != nil && current.Status == JobInvocationStatusRunning {
		current.Cancel()
	} else {
		logger.MaybeDebugfContext(ctx, js.Log, "cannot cancel; job is not runnning")
	}
	return nil
}

// RunLoop is the main scheduler loop.
// This call blocks.
// It alarms on the next runtime and forks a new routine to run the job.
// It can be aborted with the scheduler's async.Latch, or calling `.Stop()`.
// If this function exits for any reason, it will mark the scheduler as stopped.
func (js *JobScheduler) RunLoop() {
	js.Latch.Started()
	defer func() {
		js.Latch.Stopped()
		js.Latch.Reset()
	}()

	if js.JobSchedule != nil {
		js.NextRuntime = js.JobSchedule.Next(js.NextRuntime)
	}

	// if the schedule returns a zero timestamp
	// it should be interpretted as *not* to automatically
	// schedule the job to be run.
	// The run loop will return and the job scheduler will be interpretted as stopped.
	if js.NextRuntime.IsZero() {
		return
	}

	for {
		if js.NextRuntime.IsZero() {
			return
		}

		runAt := time.After(js.NextRuntime.UTC().Sub(Now()))
		select {
		case <-runAt:
			if js.CanBeScheduled() {
				if _, _, err := js.RunAsyncContext(js.Background()); err != nil {
					_ = js.error(js.Background(), err)
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

// RunAsync starts a job invocation with the BaseContext the root context.
func (js *JobScheduler) RunAsync() (*JobInvocation, <-chan struct{}, error) {
	return js.RunAsyncContext(js.Background())
}

// RunAsyncContext starts a job invocation with a given context.
func (js *JobScheduler) RunAsyncContext(ctx context.Context) (*JobInvocation, <-chan struct{}, error) {
	if !js.IsIdle() {
		return nil, nil, ex.New(ErrJobAlreadyRunning, ex.OptMessagef("job: %s", js.Name()))
	}

	ctx = js.withBaseContext(ctx)
	ctx, ji := js.withInvocationContext(ctx)
	done := make(chan struct{})
	js.SetCurrent(ji)

	var err error
	var tracer TraceFinisher
	go func() {
		defer func() {
			switch {
			case err != nil && IsJobCanceled(err):
				js.onJobCompleteCanceled(ctx) // the job was canceled, either manually or by a timeout
			case err != nil:
				js.onJobCompleteError(ctx, err) // the job completed with an error
			default:
				js.onJobCompleteSuccess(ctx) // the job completed without error
			}

			if tracer != nil {
				tracer.Finish(ctx, err) // call the trace finisher if one was started
			}
			ji.Cancel() // if the job was created with a timeout, end the timeout

			close(done)              // signal callers the job is done
			js.assignCurrentToLast() // rotate in the current to the last result
		}()

		if js.Tracer != nil {
			ctx, tracer = js.Tracer.Start(ctx, js.Name())
		}
		js.onJobBegin(ctx) // signal the job is starting

		select {
		case <-ctx.Done(): // if the timeout or cancel is triggered
			err = ErrJobCanceled // set the error to a known error
			return
		case err = <-js.safeBackgroundExec(ctx): // run the job in a background routine and catch pancis
			return
		}
	}()
	return ji, done, nil
}

// Run forces the job to run.
// This call will block.
func (js *JobScheduler) Run() {
	_, done, err := js.RunAsync()
	if err != nil {
		return
	}
	<-done
}

// RunContext runs a job with a given context as the root context.
func (js *JobScheduler) RunContext(ctx context.Context) {
	_, done, err := js.RunAsyncContext(ctx)
	if err != nil {
		return
	}
	<-done
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
	isIdle = js.Current() == nil
	return
}

//
// utility functions
//

// Current returns the current job invocation.
func (js *JobScheduler) Current() (current *JobInvocation) {
	js.currentLock.Lock()
	if js.current != nil {
		current = js.current.Clone()
	}
	js.currentLock.Unlock()
	return
}

// SetCurrent sets the current invocation, it is useful for tests etc.
func (js *JobScheduler) SetCurrent(ji *JobInvocation) {
	js.currentLock.Lock()
	js.current = ji
	js.currentLock.Unlock()
}

// Last returns the last job invocation.
func (js *JobScheduler) Last() (last *JobInvocation) {
	js.lastLock.Lock()
	if js.last != nil {
		last = js.last
	}
	js.lastLock.Unlock()
	return
}

// SetLast sets the last invocation, it is useful for tests etc.
func (js *JobScheduler) SetLast(ji *JobInvocation) {
	js.lastLock.Lock()
	js.last = ji
	js.lastLock.Unlock()
}

func (js *JobScheduler) assignCurrentToLast() {
	js.lastLock.Lock()
	js.currentLock.Lock()
	js.last = js.current
	js.current = nil
	js.currentLock.Unlock()
	js.lastLock.Unlock()
}

func (js *JobScheduler) waitCurrentComplete(ctx context.Context) {
	deadlinePoll := time.NewTicker(100 * time.Millisecond)
	defer deadlinePoll.Stop()
	for {
		if js.Current().Status != JobInvocationStatusRunning {
			return
		}
		select {
		case <-ctx.Done(): // once the timeout triggers
			return
		case <-deadlinePoll.C:
			// tick over the loop to check if the current job is complete
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

func (js *JobScheduler) withBaseContext(ctx context.Context) context.Context {
	if typed, ok := js.Job.(BackgroundProvider); ok {
		ctx = typed.Background(ctx)
	}
	ctx = logger.WithPathAppend(ctx, js.Name())
	ctx = WithJobScheduler(ctx, js)
	return ctx
}

func (js *JobScheduler) withTimeoutOrCancel(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout > 0 {
		return context.WithTimeout(ctx, timeout)
	}
	return context.WithCancel(ctx)
}

func (js *JobScheduler) withInvocationContext(ctx context.Context) (context.Context, *JobInvocation) {
	ji := NewJobInvocation(js.Name())
	ji.Parameters = MergeJobParameterValues(js.Config().ParameterValues, GetJobParameterValues(ctx))
	ctx = logger.WithPathAppend(ctx, ji.ID)
	ctx, ji.Cancel = js.withTimeoutOrCancel(ctx, js.Config().TimeoutOrDefault())
	ctx = WithJobInvocation(ctx, ji)
	ctx = WithJobParameterValues(ctx, ji.Parameters)
	return ctx, ji
}

// job lifecycle hooks

func (js *JobScheduler) onJobBegin(ctx context.Context) {
	js.currentLock.Lock()
	js.current.Started = time.Now().UTC()
	js.current.Status = JobInvocationStatusRunning
	id := js.current.ID
	js.currentLock.Unlock()

	if lifecycle := js.Lifecycle(); lifecycle.OnBegin != nil {
		lifecycle.OnBegin(ctx)
	}
	if js.Log != nil && !js.Config().SkipLoggerTrigger {
		js.logTrigger(ctx, NewEvent(FlagBegin, js.Name(), OptEventJobInvocation(id)))
	}
}

func (js *JobScheduler) onJobCompleteCanceled(ctx context.Context) {
	js.currentLock.Lock()
	js.current.Complete = time.Now().UTC()
	js.current.Status = JobInvocationStatusCanceled
	id := js.current.ID
	elapsed := js.current.Elapsed()
	js.currentLock.Unlock()

	lifecycle := js.Lifecycle()
	if lifecycle.OnCancellation != nil {
		lifecycle.OnCancellation(ctx)
	}
	if js.Log != nil && !js.Config().SkipLoggerTrigger {
		js.logTrigger(ctx, NewEvent(FlagCanceled, js.Name(), OptEventJobInvocation(id), OptEventElapsed(elapsed)))
		js.logTrigger(ctx, NewEvent(FlagComplete, js.Name(), OptEventJobInvocation(id), OptEventElapsed(elapsed)))
	}
	if lifecycle.OnComplete != nil {
		lifecycle.OnComplete(ctx)
	}
}

func (js *JobScheduler) onJobCompleteSuccess(ctx context.Context) {
	js.currentLock.Lock()
	js.current.Complete = time.Now().UTC()
	js.current.Status = JobInvocationStatusSuccess
	id := js.current.ID
	elapsed := js.current.Elapsed()
	js.currentLock.Unlock()

	lifecycle := js.Lifecycle()
	if lifecycle.OnSuccess != nil {
		lifecycle.OnSuccess(ctx)
	}
	if js.Log != nil && !js.Config().SkipLoggerTrigger {
		js.logTrigger(ctx, NewEvent(FlagSuccess, js.Name(), OptEventJobInvocation(id), OptEventElapsed(elapsed)))
		js.logTrigger(ctx, NewEvent(FlagComplete, js.Name(), OptEventJobInvocation(id), OptEventElapsed(elapsed)))
	}
	if last := js.Last(); last != nil && last.Status == JobInvocationStatusErrored {
		if lifecycle.OnFixed != nil {
			lifecycle.OnFixed(ctx)
		}
		if js.Log != nil && !js.Config().SkipLoggerTrigger {
			js.logTrigger(ctx, NewEvent(FlagFixed, js.Name(), OptEventJobInvocation(id), OptEventElapsed(elapsed)))
		}
	}
	if lifecycle.OnComplete != nil {
		lifecycle.OnComplete(ctx)
	}
}

func (js *JobScheduler) onJobCompleteError(ctx context.Context, err error) {
	js.currentLock.Lock()
	js.current.Complete = time.Now().UTC()
	js.current.Status = JobInvocationStatusErrored
	js.current.Err = err
	id := js.current.ID
	elapsed := js.current.Elapsed()
	js.currentLock.Unlock()

	//
	// error
	//

	// always log the error
	_ = js.error(ctx, err)
	lifecycle := js.Lifecycle()
	if lifecycle.OnError != nil {
		lifecycle.OnError(ctx)
	}
	if js.Log != nil && !js.Config().SkipLoggerTrigger {
		js.logTrigger(ctx, NewEvent(FlagErrored, js.Name(),
			OptEventJobInvocation(id),
			OptEventErr(err),
			OptEventElapsed(elapsed),
		))
		js.logTrigger(ctx, NewEvent(FlagComplete, js.Name(), OptEventJobInvocation(id), OptEventElapsed(elapsed)))
	}

	//
	// broken; assumes that last is set, and last was a success
	//

	if last := js.Last(); last != nil && last.Status != JobInvocationStatusErrored {
		if lifecycle.OnBroken != nil {
			lifecycle.OnBroken(ctx)
		}
		if js.Log != nil && !js.Config().SkipLoggerTrigger {
			js.logTrigger(ctx, NewEvent(FlagBroken, js.Name(),
				OptEventJobInvocation(id),
				OptEventErr(err),
				OptEventElapsed(elapsed)),
			)
		}
	}
	if lifecycle.OnComplete != nil {
		lifecycle.OnComplete(ctx)
	}
}

//
// logging helpers
//

func (js *JobScheduler) logTrigger(ctx context.Context, e logger.Event) {
	if !logger.IsLoggerSet(js.Log) {
		return
	}
	js.Log.TriggerContext(ctx, e)
}

func (js *JobScheduler) debugf(ctx context.Context, format string, args ...interface{}) {
	if !logger.IsLoggerSet(js.Log) {
		return
	}
	js.Log.DebugfContext(ctx, format, args...)
}

func (js *JobScheduler) error(ctx context.Context, err error) error {
	if !logger.IsLoggerSet(js.Log) {
		return err
	}
	js.Log.ErrorContext(ctx, err)
	return err
}
