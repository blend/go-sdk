package cron

import (
	"context"
	"time"
)

// Interface assertions.
var (
	_ Job                               = (*JobBuilder)(nil)
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
)

// NewJob returns a new job factory.
func NewJob(name string, action Action, options ...JobBuilderOption) *JobBuilder {
	jb := &JobBuilder{
		name:   name,
		action: action,
	}
	for _, option := range options {
		option(jb)
	}
	return jb
}

// JobBuilderOption is a job builder option.
type JobBuilderOption func(*JobBuilder)

// OptJobBuilderSchedule is a job builder sets the job builder schedule provder.
func OptJobBuilderSchedule(schedule Schedule) JobBuilderOption {
	return func(jb *JobBuilder) { jb.ScheduleProvider = func() Schedule { return schedule } }
}

// OptJobBuilderTimeout is a job builder sets the job builder timeout provder.
func OptJobBuilderTimeout(d time.Duration) JobBuilderOption {
	return func(jb *JobBuilder) { jb.TimeoutProvider = func() time.Duration { return d } }
}

// OptJobBuilderShutdownGracePeriod is a job builder sets the job builder shutdown grace period provder.
func OptJobBuilderShutdownGracePeriod(d time.Duration) JobBuilderOption {
	return func(jb *JobBuilder) { jb.ShutdownGracePeriodProvider = func() time.Duration { return d } }
}

// OptJobBuilderEnabledProvider is a job builder sets the job builder timeout provder.
func OptJobBuilderEnabledProvider(provider func() bool) JobBuilderOption {
	return func(jb *JobBuilder) { jb.EnabledProvider = provider }
}

// OptJobBuilderOnStart is a job builder option implementation.
func OptJobBuilderOnStart(handler func(*JobInvocation)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.OnStartHandler = handler }
}

// OptJobBuilderOnCancellation is a job builder option implementation.
func OptJobBuilderOnCancellation(handler func(*JobInvocation)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.OnCancellationHandler = handler }
}

// OptJobBuilderOnComplete is a job builder option implementation.
func OptJobBuilderOnComplete(handler func(*JobInvocation)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.OnCompleteHandler = handler }
}

// OptJobBuilderOnFailure is a job builder option implementation.
func OptJobBuilderOnFailure(handler func(*JobInvocation)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.OnFailureHandler = handler }
}

// OptJobBuilderOnBroken is a job builder option implementation.
func OptJobBuilderOnBroken(handler func(*JobInvocation)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.OnBrokenHandler = handler }
}

// OptJobBuilderOnFixed is a job builder option implementation.
func OptJobBuilderOnFixed(handler func(*JobInvocation)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.OnFixedHandler = handler }
}

// OptJobBuilderOnEnabled is a job builder option implementation.
func OptJobBuilderOnEnabled(handler func(context.Context)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.OnEnabledHandler = handler }
}

// OptJobBuilderOnDisabled is a job builder option implementation.
func OptJobBuilderOnDisabled(handler func(context.Context)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.OnDisabledHandler = handler }
}

// JobBuilder allows for job creation w/o a fully formed struct.
type JobBuilder struct {
	name   string
	action Action

	ScheduleProvider                  func() Schedule
	TimeoutProvider                   func() time.Duration
	ShutdownGracePeriodProvider       func() time.Duration
	EnabledProvider                   func() bool
	ShouldSkipLoggerListenersProvider func() bool
	ShouldSkipLoggerOutputProvider    func() bool

	OnStartHandler        func(*JobInvocation)
	OnCancellationHandler func(*JobInvocation)
	OnCompleteHandler     func(*JobInvocation)
	OnFailureHandler      func(*JobInvocation)
	OnBrokenHandler       func(*JobInvocation)
	OnFixedHandler        func(*JobInvocation)
	OnEnabledHandler      func(context.Context)
	OnDisabledHandler     func(context.Context)
}

//
// implementations of interface methods
//

// Name returns the job name.
func (jb *JobBuilder) Name() string {
	return jb.name
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
	return 0
}

// ShutdownGracePeriod returns the shutdown grace period.
func (jb *JobBuilder) ShutdownGracePeriod() time.Duration {
	if jb.ShutdownGracePeriodProvider != nil {
		return jb.ShutdownGracePeriodProvider()
	}
	return 0
}

// Enabled returns if the job is enabled.
func (jb *JobBuilder) Enabled() bool {
	if jb.EnabledProvider != nil {
		return jb.EnabledProvider()
	}
	return true
}

// ShouldSkipLoggerOutput implements the should write output provider.
func (jb *JobBuilder) ShouldSkipLoggerOutput() bool {
	if jb.ShouldSkipLoggerOutputProvider != nil {
		return jb.ShouldSkipLoggerOutputProvider()
	}
	return false
}

// ShouldSkipLoggerListeners implements the should trigger listeners provider.
func (jb *JobBuilder) ShouldSkipLoggerListeners() bool {
	if jb.ShouldSkipLoggerListenersProvider != nil {
		return jb.ShouldSkipLoggerListenersProvider()
	}
	return false
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

// Execute runs the job action if it's set.
func (jb *JobBuilder) Execute(ctx context.Context) error {
	if jb.action != nil {
		return jb.action(ctx)
	}
	return nil
}
