package cron

import (
	"context"
	"time"
)

/*
A note on the naming conventions for the below interfaces.

MethodName[Receiver|Provider] is the general pattern.

"Receiver" indicates the function will be called by the manager.

"Proivder" indicates the function will be called and is expected to return a specific value.

They're mostly the same except
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

// TimeoutProvider is an interface that allows a task to be timed out.
type TimeoutProvider interface {
	Timeout() time.Duration
}

// ShutdownGracePeriodProvider is an interface that allows a task to be given extra time to shut down.
type ShutdownGracePeriodProvider interface {
	ShutdownGracePeriod() time.Duration
}

// StatusProvider is an interface that allows a task to report its status.
type StatusProvider interface {
	Status() string
}

// SerialProvider is an optional interface that prohibits
// a task from running if another instance of the task is currently running.
type SerialProvider interface {
	Serial() bool
}

// ShouldSkipLoggerListenersProvider is a type that enables or disables logger listeners.
type ShouldSkipLoggerListenersProvider interface {
	ShouldSkipLoggerListeners() bool
}

// ShouldSkipLoggerOutputProvider is a type that enables or disables logger output for events.
type ShouldSkipLoggerOutputProvider interface {
	ShouldSkipLoggerOutput() bool
}

// EnabledProvider is an optional interface that will allow jobs to control if they're enabled.
type EnabledProvider interface {
	Enabled() bool
}

// OnStartReceiver is an interface that allows a task to be signaled when it has started.
type OnStartReceiver interface {
	OnStart(context.Context)
}

// OnCancellationReceiver is an interface that allows a task to be signaled when it has been canceled.
type OnCancellationReceiver interface {
	OnCancellation(context.Context)
}

// OnCompleteReceiver is an interface that allows a task to be signaled when it has been completed.
type OnCompleteReceiver interface {
	OnComplete(context.Context)
}

// OnFailureReceiver is an interface that allows a task to be signaled when it has been completed.
type OnFailureReceiver interface {
	OnFailure(context.Context)
}

// OnBrokenReceiver is an interface that allows a job to be signaled when it is a failure that followed
// a previous success.
type OnBrokenReceiver interface {
	OnBroken(context.Context)
}

// OnFixedReceiver is an interface that allows a jbo to be signaled when is a success that followed
// a previous failure.
type OnFixedReceiver interface {
	OnFixed(context.Context)
}

// OnDisabledReceiver is a lifecycle hook for disabled events.
type OnDisabledReceiver interface {
	OnDisabled(context.Context)
}

// OnEnabledReceiver is a lifecycle hook for enabled events.
type OnEnabledReceiver interface {
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

// HistoryProvider is a job that can persist and restore its invocation history.
type HistoryProvider interface {
	RestoreHistory(context.Context) ([]JobInvocation, error)
	PersistHistory(context.Context, []JobInvocation) error
}
