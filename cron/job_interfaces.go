package cron

import (
	"context"
	"time"
)

/*
A note on the naming conventions for the below interfaces.

MethodName[Receiver|Provider] is the general pattern.

"Receiver" indicates the function will be called by the manager.

"Provider" indicates the function will be called and is expected to return a specific value.
*/

// JobConfigProvider is a type that returns a job config.
type JobConfigProvider interface {
	JobConfig() JobConfig
}

// DescriptionProvider is a type that provides a description.
type DescriptionProvider interface {
	Description() string
}

// LabelsProvider is a type that provides labels.
type LabelsProvider interface {
	Labels() map[string]string
}

// ScheduleProvider is a type that provides a schedule for the job.
// If a job does not implement this method, it is treated as
// "OnDemand" or a job that must be triggered explicitly.
type ScheduleProvider interface {
	Schedule() Schedule
}

// OnLoadHandler is a job that can run an on load step.
type OnLoadHandler interface {
	OnLoad() error
}

// OnUnloadHandler is a job that can run an on unload step.
type OnUnloadHandler interface {
	OnUnload() error
}

// TimeoutProvider is an interface that allows a task to be timed out.
type TimeoutProvider interface {
	Timeout() time.Duration
}

// ShutdownGracePeriodProvider is an interface that allows a task to be given extra time to shut down.
type ShutdownGracePeriodProvider interface {
	ShutdownGracePeriod() time.Duration
}

// ShouldSkipLoggerListenersProvider is a type that enables or disables logger listeners.
type ShouldSkipLoggerListenersProvider interface {
	ShouldSkipLoggerListeners() bool
}

// ShouldSkipLoggerOutputProvider is a type that enables or disables logger output for events.
type ShouldSkipLoggerOutputProvider interface {
	ShouldSkipLoggerOutput() bool
}

// DisabledProvider is an optional interface that will allow jobs to control if they're disabled.
type DisabledProvider interface {
	Disabled() bool
}

// OnBeginHandler is an interface that allows a job invocation to be signaled when it has started.
type OnBeginHandler interface {
	OnBegin(context.Context)
}

// OnCancellationHandler is an interface that allows a task to be signaled when it has been canceled.
type OnCancellationHandler interface {
	OnCancellation(context.Context)
}

// OnCompleteHandler is an interface that allows a task to be signaled when it has been completed.
type OnCompleteHandler interface {
	OnComplete(context.Context)
}

// OnFailureHandler is an interface that allows a task to be signaled when it has been completed.
type OnFailureHandler interface {
	OnFailure(context.Context)
}

// OnBrokenHandler is an interface that allows a job to be signaled when it is a failure that followed
// a previous success.
type OnBrokenHandler interface {
	OnBroken(context.Context)
}

// OnFixedHandler is an interface that allows a jbo to be signaled when is a success that followed
// a previous failure.
type OnFixedHandler interface {
	OnFixed(context.Context)
}

// OnDisabledHandler is a lifecycle hook for disabled events.
type OnDisabledHandler interface {
	OnDisabled(context.Context)
}

// OnEnabledHandler is a lifecycle hook for enabled events.
type OnEnabledHandler interface {
	OnEnabled(context.Context)
}

// HistoryEnabledProvider is an optional interface that will allow jobs to control if it should track history.
type HistoryEnabledProvider interface {
	HistoryEnabled() bool
}

// HistoryMaxCountProvider is an optional interface that will allow jobs to control how many history items are tracked.
type HistoryMaxCountProvider interface {
	HistoryMaxCount() int
}

// HistoryMaxAgeProvider is an optional interface that will allow jobs to control how long to track history for.
type HistoryMaxAgeProvider interface {
	HistoryMaxAge() time.Duration
}

// HistoryPersistenceEnabledProvider is an optional interface that will allow jobs to control if it should persist history.
type HistoryPersistenceEnabledProvider interface {
	HistoryPersistenceEnabled() bool
}

// HistoryProvider is a job that can persist and restore its invocation history.
type HistoryProvider interface {
	RestoreHistory(context.Context) ([]JobInvocation, error)
	PersistHistory(context.Context, []JobInvocation) error
}
