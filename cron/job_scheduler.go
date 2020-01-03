package cron

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/mathutil"
	"github.com/blend/go-sdk/ref"
	"github.com/blend/go-sdk/stringutil"
)

// NewJobScheduler returns a job scheduler for a given job.
func NewJobScheduler(job Job, options ...JobSchedulerOption) *JobScheduler {
	js := &JobScheduler{
		latch: async.NewLatch(),
		Job:   job,
	}
	if typed, ok := job.(JobConfigProvider); ok {
		js.Config = typed.JobConfig()
	}
	for _, option := range options {
		option(js)
	}

	return js
}

// JobScheduler is a job instance.
type JobScheduler struct {
	sync.Mutex

	Job    Job
	Config JobConfig
	Tracer Tracer
	Log    logger.Log

	NextRuntime time.Time
	Current     *JobInvocation
	Last        *JobInvocation
	History     []JobInvocation

	latch    *async.Latch
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

// Description returns the description.
func (js *JobScheduler) Description() string {
	if typed, ok := js.Job.(DescriptionProvider); ok {
		return typed.Description()
	}
	return js.Config.Description
}

// Disabled returns if the job is disabled or not.
func (js *JobScheduler) Disabled() bool {
	if js.disabled != nil {
		return *js.disabled
	}
	if typed, ok := js.Job.(DisabledProvider); ok {
		return typed.Disabled()
	}
	return js.Config.DisabledOrDefault()
}

// Timeout returns the timeout or a default.
func (js *JobScheduler) Timeout() time.Duration {
	if typed, ok := js.Job.(TimeoutProvider); ok {
		return typed.Timeout()
	}
	return js.Config.TimeoutOrDefault()
}

// ShutdownGracePeriod returns the job cancellation or stop grace period.
func (js *JobScheduler) ShutdownGracePeriod() time.Duration {
	if typed, ok := js.Job.(ShutdownGracePeriodProvider); ok {
		return typed.ShutdownGracePeriod()
	}
	return js.Config.ShutdownGracePeriodOrDefault()
}

// HistoryEnabled returns if the job should track history.
func (js *JobScheduler) HistoryEnabled() bool {
	if typed, ok := js.Job.(HistoryEnabledProvider); ok {
		return typed.HistoryEnabled()
	}
	return js.Config.HistoryEnabledOrDefault()
}

// HistoryPersistenceEnabled returns if the job should call the job persistence handlers.
func (js *JobScheduler) HistoryPersistenceEnabled() bool {
	if typed, ok := js.Job.(HistoryPersistenceEnabledProvider); ok {
		return typed.HistoryPersistenceEnabled()
	}
	return js.Config.HistoryPersistenceEnabledOrDefault()
}

// HistoryMaxCount returns the maximum number of history items to keep in memory.
// 0 disables constraining by a max count.
func (js *JobScheduler) HistoryMaxCount() int {
	if typed, ok := js.Job.(HistoryMaxCountProvider); ok {
		return typed.HistoryMaxCount()
	}
	return js.Config.HistoryMaxCountOrDefault()
}

// HistoryMaxAge returns the maximum age of history items to keep in memory.
// 0 disables constraining by a max age.
func (js *JobScheduler) HistoryMaxAge() time.Duration {
	if typed, ok := js.Job.(HistoryMaxAgeProvider); ok {
		return typed.HistoryMaxAge()
	}
	return js.Config.HistoryMaxAgeOrDefault()
}

// ShouldSkipLoggerListeners returns if we should skip firing logger listeners.
func (js *JobScheduler) ShouldSkipLoggerListeners() bool {
	if typed, ok := js.Job.(ShouldSkipLoggerListenersProvider); ok {
		return typed.ShouldSkipLoggerListeners()
	}
	return js.Config.ShouldSkipLoggerListenersOrDefault()
}

// ShouldSkipLoggerOutput returns if we should have logger events skip writing to output.
// This is useful for when logging is enabled but specific jobs execute a lot.
func (js *JobScheduler) ShouldSkipLoggerOutput() bool {
	if typed, ok := js.Job.(ShouldSkipLoggerOutputProvider); ok {
		return typed.ShouldSkipLoggerOutput()
	}
	return js.Config.ShouldSkipLoggerOutputOrDefault()
}

// Labels returns the job labels, including
// automatically added ones like `name`.
func (js *JobScheduler) Labels() map[string]string {
	output := map[string]string{
		"name":      stringutil.Slugify(js.Name()),
		"scheduler": string(js.State()),
		"active":    fmt.Sprint(!js.Idle()),
		"enabled":   fmt.Sprint(!js.Disabled()),
	}
	if js.Last != nil {
		output["last"] = stringutil.Slugify(string(js.Last.State))
	}
	// config labels
	for key, value := range js.Config.Labels {
		output[key] = value
	}
	if typed, ok := js.Job.(LabelsProvider); ok {
		for key, value := range typed.Labels() {
			output[key] = stringutil.Slugify(value)
		}
	}
	return output
}

// State returns the job scheduler state.
func (js *JobScheduler) State() JobSchedulerState {
	if js.latch.IsStarted() {
		return JobSchedulerStateRunning
	}
	if js.latch.IsStopped() {
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
		Timeout:                   js.Timeout(),
		NextRuntime:               js.NextRuntime,
		Current:                   js.Current,
		Last:                      js.Last,
		Stats:                     js.Stats(),
		HistoryEnabled:            js.HistoryEnabled(),
		HistoryPersistenceEnabled: js.HistoryPersistenceEnabled(),
		HistoryMaxCount:           js.HistoryMaxCount(),
		HistoryMaxAge:             js.HistoryMaxAge(),
	}
	if js.Schedule() != nil {
		if typed, ok := js.Schedule().(fmt.Stringer); ok {
			status.Schedule = typed.String()
		}
	}

	status.History = make([]JobSchedulerStatusHistory, len(js.History))
	for index, ji := range js.History {
		status.History[index] = JobSchedulerStatusHistory{
			Started: ji.Started,
			State:   ji.State,
			Elapsed: ji.Elapsed,
		}
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
	if !js.latch.CanStart() {
		return fmt.Errorf("already started")
	}
	js.infof("scheduler starting")
	js.latch.Starting()
	js.infof("scheduler started")
	js.RunLoop()
	js.infof("scheduler exiting")
	return nil
}

// Stop stops the scheduler.
func (js *JobScheduler) Stop() error {
	if !js.latch.CanStop() {
		return fmt.Errorf("already stopped")
	}
	stopped := js.latch.NotifyStopped()
	js.infof("scheduler stopping")
	// signal we are stopping.
	js.latch.Stopping()

	if js.Current != nil {
		gracePeriod := js.ShutdownGracePeriod()
		if gracePeriod > 0 {
			js.debugf("job cancellation; cancelling with %v grace period", gracePeriod)
			ctx, cancel := js.createContextWithTimeout(context.Background(), gracePeriod)
			defer cancel()

			js.cancelJobInvocation(ctx, js.Current)
		} else {
			js.cancelJobInvocation(context.Background(), js.Current)
		}
	}
	js.PersistHistory(context.Background())

	<-stopped
	js.latch.Reset()
	js.infof("scheduler stopped")
	return nil
}

// NotifyStarted notifies the job scheduler has started.
func (js *JobScheduler) NotifyStarted() <-chan struct{} {
	return js.latch.NotifyStarted()
}

// NotifyStopped notifies the job scheduler has stopped.
func (js *JobScheduler) NotifyStopped() <-chan struct{} {
	return js.latch.NotifyStopped()
}

// Enable sets the job as enabled.
func (js *JobScheduler) Enable() {
	js.disabled = ref.Bool(false)
	if js.Log != nil && !js.ShouldSkipLoggerListeners() {
		js.Log.Trigger(js.logEventContext(context.Background()), NewEvent(FlagEnabled, js.Name()))
	}
	if typed, ok := js.Job.(OnEnabledHandler); ok {
		typed.OnEnabled(context.Background())
	}
}

// Disable sets the job as disabled.
func (js *JobScheduler) Disable() {
	js.disabled = ref.Bool(true)
	if js.Log != nil && !js.ShouldSkipLoggerListeners() {
		js.Log.Trigger(js.logEventContext(context.Background()), NewEvent(FlagDisabled, js.Name()))
	}
	if typed, ok := js.Job.(OnDisabledHandler); ok {
		typed.OnDisabled(context.Background())
	}
}

// Cancel stops all running invocations.
func (js *JobScheduler) Cancel() error {
	if js.Current == nil {
		js.debugf("job cancellation; not running")
		return nil
	}
	gracePeriod := js.ShutdownGracePeriod()
	if gracePeriod > 0 {
		js.debugf("job cancellation; cancelling with %v grace period", gracePeriod)
		ctx, cancel := js.createContextWithTimeout(context.Background(), gracePeriod)
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
	js.latch.Started()
	defer func() {
		js.latch.Stopped()
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
	notifyStopping := js.latch.NotifyStopping()

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
	if !js.Idle() {
		return nil, ex.New(ErrJobAlreadyRunning, ex.OptMessagef("job: %s", js.Name()))
	}

	timeout := js.Timeout()

	// create a job invocation, or a record of each
	// individual execution of a job.
	ji := NewJobInvocation(js.Name())
	ji.Context, ji.Cancel = js.createContextWithTimeout(ctx, timeout)
	ji.Parameters = GetJobParameters(ctx) // pull the parameters off the calling context.

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
		// fire the on begin event
		js.onBegin(ji.Context, ji)

		// check if the job has been canceled
		// or if it's finished.
		select {
		case <-ji.Context.Done():
			err = ErrJobCancelled
			return
		case err = <-js.safeBackgroundExec(ji.Context): // we use a goroutine here
			return // so that we can enforce the timeout
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
	if !js.HistoryPersistenceEnabled() {
		return nil
	}

	historyProvider, ok := js.Job.(HistoryProvider)
	if !ok {
		return nil
	}

	js.Lock()
	defer js.Unlock()
	var err error
	if js.History, err = historyProvider.RestoreHistory(ctx); err != nil {
		return js.error(err)
	}
	if len(js.History) > 0 {
		js.Last = &js.History[len(js.History)-1]
	}
	return nil
}

// PersistHistory calls the persist handler if it's set.
func (js *JobScheduler) PersistHistory(ctx context.Context) error {
	if !js.HistoryEnabled() {
		return nil
	}
	if !js.HistoryPersistenceEnabled() {
		return nil
	}

	historyProvider, ok := js.Job.(HistoryProvider)
	if !ok {
		return nil
	}
	js.Lock()
	defer js.Unlock()

	historyCopy := make([]JobInvocation, len(js.History))
	copy(historyCopy, js.History)
	if err := historyProvider.PersistHistory(ctx, historyCopy); err != nil {
		return js.error(err)
	}
	return nil
}

// CanBeScheduled returns if a job will be triggered automatically
// and isn't already in flight and set to be serial.
func (js *JobScheduler) CanBeScheduled() bool {
	return !js.Disabled() && js.Idle()
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

	if js.HistoryEnabled() {
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
	maxCount := js.HistoryMaxCount()
	maxAge := js.HistoryMaxAge()

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

func (js *JobScheduler) onBegin(ctx context.Context, ji *JobInvocation) {
	js.Lock()
	defer js.Unlock()

	if js.Log != nil && !js.ShouldSkipLoggerListeners() {
		js.logTrigger(js.logEventContext(ctx), NewEvent(FlagBegin, ji.JobName, OptEventJobInvocation(ji.ID)))
	}
	if typed, ok := js.Job.(OnBeginHandler); ok {
		typed.OnBegin(ctx)
	}
}

func (js *JobScheduler) onCancelled(ctx context.Context, ji *JobInvocation) {
	js.Lock()
	defer js.Unlock()

	ji.State = JobInvocationStateCancelled

	if js.Log != nil && !js.ShouldSkipLoggerListeners() {
		js.logTrigger(js.logEventContext(ctx), NewEvent(FlagCancelled, ji.JobName, OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed)))
	}
	if typed, ok := js.Job.(OnCancellationHandler); ok {
		typed.OnCancellation(ctx)
	}
}

func (js *JobScheduler) onComplete(ctx context.Context, ji *JobInvocation) {
	js.Lock()
	defer js.Unlock()

	ji.State = JobInvocationStateComplete

	if js.Log != nil && !js.ShouldSkipLoggerListeners() {
		js.logTrigger(js.logEventContext(ctx), NewEvent(FlagComplete, ji.JobName, OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed)))
	}
	if typed, ok := js.Job.(OnCompleteHandler); ok {
		typed.OnComplete(ctx)
	}

	if js.Last != nil && js.Last.Err != nil {
		js.logTrigger(js.logEventContext(ctx), NewEvent(FlagFixed, ji.JobName, OptEventElapsed(ji.Elapsed)))
		if typed, ok := js.Job.(OnFixedHandler); ok {
			typed.OnFixed(ctx)
		}
	}
}

func (js *JobScheduler) onFailure(ctx context.Context, ji *JobInvocation) {
	js.Lock()
	defer js.Unlock()

	ji.State = JobInvocationStateFailed

	if js.Log != nil && !js.ShouldSkipLoggerListeners() {
		js.logTrigger(js.logEventContext(ctx), NewEvent(FlagFailed, ji.JobName, OptEventErr(ji.Err), OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed)))
	}
	if ji.Err != nil {
		js.error(ji.Err)
	}
	if typed, ok := js.Job.(OnFailureHandler); ok {
		typed.OnFailure(ctx)
	}
	if js.Last != nil && js.Last.Err == nil {
		if js.Log != nil {
			js.logTrigger(js.logEventContext(ctx), NewEvent(FlagBroken, ji.JobName, OptEventJobInvocation(ji.ID), OptEventElapsed(ji.Elapsed)))
		}
		if typed, ok := js.Job.(OnBrokenHandler); ok {
			typed.OnBroken(ctx)
		}
	}
}

//
// logging helpers
//

func (js *JobScheduler) logEventContext(parent context.Context) context.Context {
	if js.ShouldSkipLoggerOutput() {
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
