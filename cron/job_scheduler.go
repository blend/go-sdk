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
	Tracer Tracer
	Log    logger.Log

	NextRuntime time.Time
	Current     *JobInvocation
	Last        *JobInvocation
	History     []JobInvocation

	disabled *bool
}

// Name returns the job name.
func (js *JobScheduler) Name() string {
	return js.Job.Name()
}

// Schedule returns the job schedule.
func (js *JobScheduler) Schedule() Schedule {
	if typed, ok := js.Job.(ScheduleProvider); ok {
		return typed.Schedule()
	}
	return nil
}

// Config returns the job config provided by a job or an empty config.
func (js *JobScheduler) Config() JobConfig {
	if typed, ok := js.Job.(ConfigProvider); ok {
		return typed.Config()
	}
	return JobConfig{}
}

// Lifecycle returns job lifecycle steps or an empty set.
func (js *JobScheduler) Lifecycle() JobLifecycle {
	if typed, ok := js.Job.(LifecycleProvider); ok {
		return typed.Lifecycle()
	}
	return JobLifecycle{}
}

// Description returns the description.
func (js *JobScheduler) Description() string {
	return js.Config().Description
}

// Disabled returns if the job is disabled or not.
func (js *JobScheduler) Disabled() bool {
	if js.disabled != nil {
		return *js.disabled
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

// Status returns the job scheduler status.
func (js *JobScheduler) Status() JobSchedulerStatus {
	status := JobSchedulerStatus{
		Name:                      js.Name(),
		State:                     js.State(),
		Labels:                    js.Labels(),
		Disabled:                  js.Disabled(),
		Timeout:                   js.Config().TimeoutOrDefault(),
		NextRuntime:               js.NextRuntime,
		Current:                   js.Current,
		Last:                      js.Last,
		HistoryEnabled:            js.Config().HistoryEnabledOrDefault(),
		HistoryPersistenceEnabled: js.Config().HistoryPersistenceEnabledOrDefault(),
		HistoryMaxCount:           js.Config().HistoryMaxCountOrDefault(),
		HistoryMaxAge:             js.Config().HistoryMaxAgeOrDefault(),
	}
	if js.Schedule() != nil {
		if typed, ok := js.Schedule().(fmt.Stringer); ok {
			status.Schedule = typed.String()
		}
	}
	return status
}

// Start starts the scheduler.
// This call blocks.
func (js *JobScheduler) Start() error {
	if !js.Latch.CanStart() {
		return fmt.Errorf("already started")
	}
	js.Latch.Starting()
	js.RunLoop()
	return nil
}

// Stop stops the scheduler.
func (js *JobScheduler) Stop() error {
	if !js.Latch.CanStop() {
		return fmt.Errorf("already stopped")
	}
	stopped := js.Latch.NotifyStopped()
	js.Latch.Stopping()

	if js.Current != nil {
		gracePeriod := js.Config().ShutdownGracePeriodOrDefault()
		if gracePeriod > 0 {
			ctx, cancel := js.createContextWithTimeout(context.Background(), gracePeriod)
			defer cancel()
			js.cancelJobInvocation(ctx, js.Current)
		} else {
			js.cancelJobInvocation(context.Background(), js.Current)
		}
	}
	js.PersistHistory(context.Background())
	<-stopped
	js.Latch.Reset()
	return nil
}

// OnLoad triggers the on load even on the job lifecycle handler.
func (js *JobScheduler) OnLoad() error {
	if js.Lifecycle().OnLoad != nil {
		return js.Lifecycle().OnLoad()
	}
	return nil
}

// OnUnload triggers the on unload even on the job lifecycle handler.
func (js *JobScheduler) OnUnload() error {
	if js.Lifecycle().OnUnload != nil {
		return js.Lifecycle().OnUnload()
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
	js.disabled = ref.Bool(false)
	if lifecycle := js.Lifecycle(); lifecycle.OnEnabled != nil {
		lifecycle.OnEnabled(js.logEventContext(context.Background(), nil))
	}
	if js.Log != nil && !js.Config().ShouldSkipLoggerListenersOrDefault() {
		js.Log.Trigger(js.logEventContext(context.Background(), nil), NewEvent(FlagEnabled, js.Name()))
	}
}

// Disable sets the job as disabled.
func (js *JobScheduler) Disable() {
	js.disabled = ref.Bool(true)
	if lifecycle := js.Lifecycle(); lifecycle.OnDisabled != nil {
		lifecycle.OnDisabled(js.logEventContext(context.Background(), nil))
	}
	if js.Log != nil && !js.Config().ShouldSkipLoggerListenersOrDefault() {
		js.Log.Trigger(js.logEventContext(context.Background(), nil), NewEvent(FlagDisabled, js.Name()))
	}
}

// Cancel stops all running invocations.
func (js *JobScheduler) Cancel() error {
	if js.Current == nil {
		return nil
	}
	gracePeriod := js.Config().ShutdownGracePeriodOrDefault()
	if gracePeriod > 0 {
		ctx, cancel := js.createContextWithTimeout(context.Background(), gracePeriod)
		defer cancel()

		js.cancelJobInvocation(ctx, js.Current)
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
	}()

	if js.Schedule() != nil {
		js.NextRuntime = js.Schedule().Next(js.NextRuntime)
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
			if js.Schedule() != nil {
				js.NextRuntime = js.Schedule().Next(js.NextRuntime)
			} else {
				js.NextRuntime = time.Time{}
			}

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
	if !js.IsIdle() {
		return nil, ex.New(ErrJobAlreadyRunning, ex.OptMessagef("job: %s", js.Name()))
	}

	timeout := js.Config().TimeoutOrDefault()

	// create a job invocation, or a record of each
	// individual execution of a job.
	ji := NewJobInvocation(js.Name())
	ctx = js.logEventContext(ctx, ji)
	ji.Context, ji.Cancel = js.createContextWithTimeout(ctx, timeout)
	ji.Parameters = GetJobParameters(ctx) // pull the parameters off the calling context.

	if timeout > 0 {
		ji.Timeout = ji.Started.Add(timeout)
	}
	js.setCurrent(ji)

	var err error
	var tracer TraceFinisher
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
			if tracer != nil {
				tracer.Finish(ji.Context, err)
			}

			if err != nil && IsJobCancelled(err) {
				js.onJobCancelled(ji.Context, ji)
			} else if err != nil {
				ji.Err = err
				js.onJobError(ji.Context, ji)
			} else {
				js.onJobComplete(ji.Context, ji)
			}
			js.setLast(ji)
			js.PersistHistory(ji.Context)
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

// GetJobInvocationByID returns an invocation by id.
func (js *JobScheduler) GetJobInvocationByID(id string) *JobInvocation {
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

// History functions

// RestoreHistory calls the persist handler if it's set.
func (js *JobScheduler) RestoreHistory(ctx context.Context) error {
	if !js.Config().HistoryPersistenceEnabledOrDefault() {
		return nil
	}
	if js.Lifecycle().RestoreHistory == nil {
		return nil
	}

	js.Lock()
	defer js.Unlock()
	var err error
	if js.History, err = js.Lifecycle().RestoreHistory(ctx); err != nil {
		return js.error(ctx, err)
	}
	if len(js.History) > 0 {
		js.Last = &js.History[len(js.History)-1]
	}
	return nil
}

// PersistHistory calls the persist handler if it's set.
func (js *JobScheduler) PersistHistory(ctx context.Context) error {
	if !js.Config().HistoryEnabledOrDefault() {
		return nil
	}
	if !js.Config().HistoryPersistenceEnabledOrDefault() {
		return nil
	}
	if js.Lifecycle().PersistHistory == nil {
		return nil
	}

	js.Lock()
	defer js.Unlock()

	historyCopy := make([]JobInvocation, len(js.History))
	copy(historyCopy, js.History)
	if err := js.Lifecycle().PersistHistory(ctx, historyCopy); err != nil {
		return js.error(ctx, err)
	}
	return nil
}

//
// utility functions
//

func (js *JobScheduler) setLast(ji *JobInvocation) {
	js.Lock()
	defer js.Unlock()

	if js.Config().HistoryEnabledOrDefault() {
		js.History = append(js.cullHistory(), *ji)
	}
	js.Current = nil
	js.Last = ji
	close(ji.Done)
}

func (js *JobScheduler) setCurrent(ji *JobInvocation) {
	js.Lock()
	js.Current = ji
	js.Unlock()
}

func (js *JobScheduler) cullHistory() []JobInvocation {
	count := len(js.History)
	maxCount := js.Config().HistoryMaxCountOrDefault()
	maxAge := js.Config().HistoryMaxAgeOrDefault()

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
		if ji == nil || ji.Status != JobInvocationStatusRunning {
			return
		}
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

// job lifecycle hooks

func (js *JobScheduler) onJobBegin(ctx context.Context, ji *JobInvocation) {
	js.Lock()
	defer js.Unlock()

	ji.Status = JobInvocationStatusRunning

	if js.Log != nil && !js.Config().ShouldSkipLoggerListenersOrDefault() {
		js.logTrigger(ctx, NewEvent(FlagBegin, ji.JobName, OptEventJobInvocation(ji.ID)))
	}
	if lifecycle := js.Lifecycle(); lifecycle.OnBegin != nil {
		lifecycle.OnBegin(ctx)
	}
}

func (js *JobScheduler) onJobCancelled(ctx context.Context, ji *JobInvocation) {
	js.Lock()
	defer js.Unlock()

	ji.Status = JobInvocationStatusCancelled

	if js.Log != nil && !js.Config().ShouldSkipLoggerListenersOrDefault() {
		js.logTrigger(ctx, NewEvent(FlagCancelled, ji.JobName, OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed())))
	}
	if lifecycle := js.Lifecycle(); lifecycle.OnCancellation != nil {
		lifecycle.OnCancellation(ctx)
	}
}

func (js *JobScheduler) onJobComplete(ctx context.Context, ji *JobInvocation) {
	js.Lock()
	defer js.Unlock()

	ji.Status = JobInvocationStatusComplete

	if js.Log != nil && !js.Config().ShouldSkipLoggerListenersOrDefault() {
		js.logTrigger(ctx, NewEvent(FlagComplete, ji.JobName, OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed())))
	}
	if lifecycle := js.Lifecycle(); lifecycle.OnComplete != nil {
		lifecycle.OnComplete(ctx)
	}

	if js.Last != nil && js.Last.Err != nil {
		if lifecycle := js.Lifecycle(); lifecycle.OnFixed != nil {
			lifecycle.OnFixed(ctx)
		}
		js.logTrigger(ctx, NewEvent(FlagFixed, ji.JobName, OptEventElapsed(ji.Elapsed())))
	}
}

func (js *JobScheduler) onJobError(ctx context.Context, ji *JobInvocation) {
	js.Lock()
	defer js.Unlock()

	ji.Status = JobInvocationStatusErrored

	if js.Log != nil && !js.Config().ShouldSkipLoggerListenersOrDefault() {
		js.logTrigger(ctx, NewEvent(FlagErrored, ji.JobName, OptEventErr(ji.Err), OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed())))
	}
	if lifecycle := js.Lifecycle(); lifecycle.OnError != nil {
		lifecycle.OnError(ctx)
	}
	if ji.Err != nil {
		js.error(ctx, ji.Err)
	}

	if js.Last != nil && js.Last.Err == nil {
		if lifecycle := js.Lifecycle(); lifecycle.OnBroken != nil {
			lifecycle.OnBroken(ctx)
		}
		if js.Log != nil {
			js.logTrigger(ctx, NewEvent(FlagBroken, ji.JobName, OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed())))
		}
	}
}

//
// logging helpers
//

func (js *JobScheduler) logEventContext(parent context.Context, ji *JobInvocation) context.Context {
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

func (js *JobScheduler) error(ctx context.Context, err error) error {
	if js.Log == nil {
		return err
	}
	js.Log.WithContext(ctx).Error(err)
	return err
}
