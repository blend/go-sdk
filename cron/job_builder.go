package cron

import (
	"context"
	"time"

	"github.com/blend/go-sdk/stringutil"
)

// Interface assertions.
var (
	_ Job                               = (*JobBuilder)(nil)
	_ LabelsProvider                    = (*JobBuilder)(nil)
	_ ScheduleProvider                  = (*JobBuilder)(nil)
	_ TimeoutProvider                   = (*JobBuilder)(nil)
	_ ShutdownGracePeriodProvider       = (*JobBuilder)(nil)
	_ EnabledProvider                   = (*JobBuilder)(nil)
	_ ShouldSkipLoggerListenersProvider = (*JobBuilder)(nil)
	_ ShouldSkipLoggerOutputProvider    = (*JobBuilder)(nil)
	_ OnStartReceiver                   = (*JobBuilder)(nil)
	_ OnCancellationReceiver            = (*JobBuilder)(nil)
	_ OnCompleteReceiver                = (*JobBuilder)(nil)
	_ OnFailureReceiver                 = (*JobBuilder)(nil)
	_ OnBrokenReceiver                  = (*JobBuilder)(nil)
	_ OnFixedReceiver                   = (*JobBuilder)(nil)
	_ OnEnabledReceiver                 = (*JobBuilder)(nil)
	_ OnDisabledReceiver                = (*JobBuilder)(nil)
	_ HistoryEnabledProvider            = (*JobBuilder)(nil)
	_ HistoryMaxCountProvider           = (*JobBuilder)(nil)
	_ HistoryMaxAgeProvider             = (*JobBuilder)(nil)
	_ HistoryProvider                   = (*JobBuilder)(nil)
)

// NewJob returns a new job builder.
func NewJob(options ...JobBuilderOption) *JobBuilder {
	var jb JobBuilder
	for _, option := range options {
		option(&jb)
	}
	if jb.Config.Name == "" {
		jb.Config.Name = stringutil.Random(stringutil.LowerLetters, 16)
	}
	return &jb
}

// JobBuilderOption is a job builder option.
type JobBuilderOption func(*JobBuilder)

// OptJobName sets the job name.
func OptJobName(name string) JobBuilderOption {
	return func(jb *JobBuilder) { jb.Config.Name = name }
}

// OptJobAction sets the job action.
func OptJobAction(action func(context.Context) error) JobBuilderOption {
	return func(jb *JobBuilder) { jb.Action = action }
}

// OptJobSchedule is a job builder sets the job schedule provder.
func OptJobSchedule(schedule Schedule) JobBuilderOption {
	return func(jb *JobBuilder) { jb.ScheduleProvider = func() Schedule { return schedule } }
}

// OptJobTimeout is a job builder sets the job timeout provder.
func OptJobTimeout(d time.Duration) JobBuilderOption {
	return func(jb *JobBuilder) { jb.TimeoutProvider = func() time.Duration { return d } }
}

// OptJobShutdownGracePeriod is a job builder sets the job shutdown grace period provder.
func OptJobShutdownGracePeriod(d time.Duration) JobBuilderOption {
	return func(jb *JobBuilder) { jb.ShutdownGracePeriodProvider = func() time.Duration { return d } }
}

// OptJobEnabled is a job builder sets the job timeout provder.
func OptJobEnabled(provider func() bool) JobBuilderOption {
	return func(jb *JobBuilder) { jb.EnabledProvider = provider }
}

// OptJobOnStart is a job builder option implementation.
func OptJobOnStart(handler func(*JobInvocation)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.OnStartHandler = handler }
}

// OptJobOnCancellation is a job builder option implementation.
func OptJobOnCancellation(handler func(*JobInvocation)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.OnCancellationHandler = handler }
}

// OptJobOnComplete is a job builder option implementation.
func OptJobOnComplete(handler func(*JobInvocation)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.OnCompleteHandler = handler }
}

// OptJobOnFailure is a job builder option implementation.
func OptJobOnFailure(handler func(*JobInvocation)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.OnFailureHandler = handler }
}

// OptJobOnBroken is a job builder option implementation.
func OptJobOnBroken(handler func(*JobInvocation)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.OnBrokenHandler = handler }
}

// OptJobOnFixed is a job builder option implementation.
func OptJobOnFixed(handler func(*JobInvocation)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.OnFixedHandler = handler }
}

// OptJobOnEnabled is a job builder option implementation.
func OptJobOnEnabled(handler func(context.Context)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.OnEnabledHandler = handler }
}

// OptJobOnDisabled is a job builder option implementation.
func OptJobOnDisabled(handler func(context.Context)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.OnDisabledHandler = handler }
}

// OptJobHistoryEnabled is a job builder option implementation.
func OptJobHistoryEnabled(provider func() bool) JobBuilderOption {
	return func(jb *JobBuilder) { jb.HistoryEnabledProvider = provider }
}

// OptJobHistoryMaxCount is a job builder option implementation.
func OptJobHistoryMaxCount(provider func() int) JobBuilderOption {
	return func(jb *JobBuilder) { jb.HistoryMaxCountProvider = provider }
}

// OptJobHistoryMaxAge is a job builder option implementation.
func OptJobHistoryMaxAge(provider func() time.Duration) JobBuilderOption {
	return func(jb *JobBuilder) { jb.HistoryMaxAgeProvider = provider }
}

// JobBuilder allows for job creation w/o a fully formed struct.
type JobBuilder struct {
	Action Action
	Config JobConfig

	ScheduleProvider                  func() Schedule
	TimeoutProvider                   func() time.Duration
	ShutdownGracePeriodProvider       func() time.Duration
	EnabledProvider                   func() bool
	ShouldSkipLoggerListenersProvider func() bool
	ShouldSkipLoggerOutputProvider    func() bool
	HistoryEnabledProvider            func() bool
	HistoryMaxCountProvider           func() int
	HistoryMaxAgeProvider             func() time.Duration

	OnStartHandler        func(*JobInvocation)
	OnCancellationHandler func(*JobInvocation)
	OnCompleteHandler     func(*JobInvocation)
	OnFailureHandler      func(*JobInvocation)
	OnBrokenHandler       func(*JobInvocation)
	OnFixedHandler        func(*JobInvocation)
	OnEnabledHandler      func(context.Context)
	OnDisabledHandler     func(context.Context)

	RestoreHistoryHandler func(context.Context) ([]JobInvocation, error)
	PersistHistoryHandler func(context.Context, []JobInvocation) error
}

//
// implementations of interface methods
//

// Name returns the job name.
func (jb *JobBuilder) Name() string {
	return jb.Config.Name
}

// Labels returns the job labels.
func (jb *JobBuilder) Labels() map[string]string {
	return jb.Config.Labels
}

// Schedule returns the job schedule if a provider is set.
func (jb *JobBuilder) Schedule() Schedule {
	if jb.ScheduleProvider != nil {
		return jb.ScheduleProvider()
	}
	return nil
}

// Timeout returns the job timeout.
func (jb *JobBuilder) Timeout() time.Duration {
	if jb.TimeoutProvider != nil {
		return jb.TimeoutProvider()
	}
	return jb.Config.TimeoutOrDefault()
}

// ShutdownGracePeriod returns the shutdown grace period.
func (jb *JobBuilder) ShutdownGracePeriod() time.Duration {
	if jb.ShutdownGracePeriodProvider != nil {
		return jb.ShutdownGracePeriodProvider()
	}
	return jb.Config.ShutdownGracePeriodOrDefault()
}

// Enabled returns if the job is enabled.
func (jb *JobBuilder) Enabled() bool {
	if jb.EnabledProvider != nil {
		return jb.EnabledProvider()
	}
	return jb.Config.EnabledOrDefault()
}

// ShouldSkipLoggerListeners implements the should skip logger listeners provider.
func (jb *JobBuilder) ShouldSkipLoggerListeners() bool {
	if jb.ShouldSkipLoggerListenersProvider != nil {
		return jb.ShouldSkipLoggerListenersProvider()
	}
	return jb.Config.ShouldSkipLoggerListenersOrDefault()
}

// ShouldSkipLoggerOutput implements the should skip logger output provider.
func (jb *JobBuilder) ShouldSkipLoggerOutput() bool {
	if jb.ShouldSkipLoggerOutputProvider != nil {
		return jb.ShouldSkipLoggerOutputProvider()
	}
	return jb.Config.ShouldSkipLoggerOutputOrDefault()
}

// HistoryEnabled implements the history enabled provider.
func (jb *JobBuilder) HistoryEnabled() bool {
	if jb.HistoryEnabledProvider != nil {
		return jb.HistoryEnabledProvider()
	}
	return jb.Config.HistoryEnabledOrDefault()
}

// HistoryMaxCount implements the history max count provider.
func (jb *JobBuilder) HistoryMaxCount() int {
	if jb.HistoryMaxCountProvider != nil {
		return jb.HistoryMaxCountProvider()
	}
	return jb.Config.HistoryMaxCountOrDefault()
}

// HistoryMaxAge implements the history max count provider.
func (jb *JobBuilder) HistoryMaxAge() time.Duration {
	if jb.HistoryMaxAgeProvider != nil {
		return jb.HistoryMaxAgeProvider()
	}
	return jb.Config.HistoryMaxAgeOrDefault()
}

// OnStart is a lifecycle hook.
func (jb *JobBuilder) OnStart(ctx context.Context) {
	if jb.OnStartHandler != nil {
		jb.OnStartHandler(GetJobInvocation(ctx))
	}
}

// OnCancellation is a lifecycle hook.
func (jb *JobBuilder) OnCancellation(ctx context.Context) {
	if jb.OnCancellationHandler != nil {
		jb.OnCancellationHandler(GetJobInvocation(ctx))
	}
}

// OnComplete is a lifecycle hook.
func (jb *JobBuilder) OnComplete(ctx context.Context) {
	if jb.OnCompleteHandler != nil {
		jb.OnCompleteHandler(GetJobInvocation(ctx))
	}
}

// OnFailure is a lifecycle hook.
func (jb *JobBuilder) OnFailure(ctx context.Context) {
	if jb.OnFailureHandler != nil {
		jb.OnFailureHandler(GetJobInvocation(ctx))
	}
}

// OnFixed is a lifecycle hook.
func (jb *JobBuilder) OnFixed(ctx context.Context) {
	if jb.OnFixedHandler != nil {
		jb.OnFixedHandler(GetJobInvocation(ctx))
	}
}

// OnBroken is a lifecycle hook.
func (jb *JobBuilder) OnBroken(ctx context.Context) {
	if jb.OnBrokenHandler != nil {
		jb.OnBrokenHandler(GetJobInvocation(ctx))
	}
}

// OnEnabled is a lifecycle hook.
func (jb *JobBuilder) OnEnabled(ctx context.Context) {
	if jb.OnEnabledHandler != nil {
		jb.OnEnabledHandler(ctx)
	}
}

// OnDisabled is a lifecycle hook.
func (jb *JobBuilder) OnDisabled(ctx context.Context) {
	if jb.OnDisabledHandler != nil {
		jb.OnDisabledHandler(ctx)
	}
}

// RestoreHistory calls the restore history handler if it's set.
func (jb *JobBuilder) RestoreHistory(ctx context.Context) ([]JobInvocation, error) {
	if jb.RestoreHistoryHandler != nil {
		return jb.RestoreHistoryHandler(ctx)
	}
	return nil, nil
}

// PersistHistory calls the persist history handler if it's set.
func (jb *JobBuilder) PersistHistory(ctx context.Context, history []JobInvocation) error {
	if jb.PersistHistoryHandler != nil {
		return jb.PersistHistoryHandler(ctx, history)
	}
	return nil
}

// Execute runs the job action if it's set.
func (jb *JobBuilder) Execute(ctx context.Context) error {
	if jb.Action != nil {
		return jb.Action(ctx)
	}
	return nil
}
