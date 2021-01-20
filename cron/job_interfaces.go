package cron

import (
	"context"
)

// Job is an interface types can satisfy to be loaded into the JobManager.
type Job interface {
	Name() string
	Execute(context.Context) error
}

// ConfigProvider is a type that returns a job config.
type ConfigProvider interface {
	Config() JobConfig
}

// BackgroundProvider is a type that returns a base context based on a parent.
type BackgroundProvider interface {
	Background(context.Context) context.Context
}

// ScheduleProvider is a type that provides a schedule for the job.
// If a job does not implement this method, it is treated as
// "OnDemand" or a job that must be triggered explicitly.
type ScheduleProvider interface {
	Schedule() Schedule
}

// LifecycleProvider is a job that provides lifecycle hooks.
type LifecycleProvider interface {
	Lifecycle() JobLifecycle
}
