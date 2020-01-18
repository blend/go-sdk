package cron

import (
	"context"
	"time"

	"github.com/blend/go-sdk/ref"
	"github.com/blend/go-sdk/stringutil"
)

// Interface assertions.
var (
	_ Job               = (*JobBuilder)(nil)
	_ ScheduleProvider  = (*JobBuilder)(nil)
	_ LifecycleProvider = (*JobBuilder)(nil)
	_ ConfigProvider    = (*JobBuilder)(nil)
)

// NewJob returns a new job builder.
func NewJob(options ...JobBuilderOption) *JobBuilder {
	var jb JobBuilder
	for _, option := range options {
		option(&jb)
	}
	if jb.JobName == "" {
		jb.JobName = stringutil.Random(stringutil.LowerLetters, 16)
	}
	return &jb
}

// JobBuilderOption is a job builder option.
type JobBuilderOption func(*JobBuilder)

// OptJobName sets the job name.
func OptJobName(name string) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobName = name }
}

// OptJobAction sets the job action.
func OptJobAction(action func(context.Context) error) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobAction = action }
}

// OptJobLabels is a job builder sets the job labels.
func OptJobLabels(labels map[string]string) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobConfig.Labels = labels }
}

// OptJobSchedule is a job builder sets the job schedule provder.
func OptJobSchedule(schedule Schedule) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobScheduleProvider = func() Schedule { return schedule } }
}

// OptJobTimeout is a job builder sets the job timeout provder.
func OptJobTimeout(d time.Duration) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobConfig.Timeout = d }
}

// OptJobShutdownGracePeriod is a job builder sets the job shutdown grace period provder.
func OptJobShutdownGracePeriod(d time.Duration) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobConfig.ShutdownGracePeriod = d }
}

// OptJobDisabled is a job builder sets the job timeout provder.
func OptJobDisabled(disabled bool) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobConfig.Disabled = ref.Bool(disabled) }
}

// OptJobHistoryEnabled is a job builder option implementation.
func OptJobHistoryEnabled(enabled bool) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobConfig.HistoryEnabled = ref.Bool(enabled) }
}

// OptJobHistoryPersistenceEnabled is a job builder option implementation.
func OptJobHistoryPersistenceEnabled(enabled bool) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobConfig.HistoryPersistenceEnabled = ref.Bool(enabled) }
}

// OptJobHistoryMaxCount is a job builder option implementation.
func OptJobHistoryMaxCount(maxCount int) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobConfig.HistoryMaxCount = ref.Int(maxCount) }
}

// OptJobHistoryMaxAge is a job builder option implementation.
func OptJobHistoryMaxAge(maxAge time.Duration) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobConfig.HistoryMaxAge = ref.Duration(maxAge) }
}

// OptJobOnBegin sets a lifecycle hook.
func OptJobOnBegin(handler func(context.Context)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobLifecycle.OnBegin = handler }
}

// OptJobOnComplete sets a lifecycle hook.
func OptJobOnComplete(handler func(context.Context)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobLifecycle.OnComplete = handler }
}

// OptJobOnError sets a lifecycle hook.
func OptJobOnError(handler func(context.Context)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobLifecycle.OnError = handler }
}

// OptJobOnCancellation sets a lifecycle hook.
func OptJobOnCancellation(handler func(context.Context)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobLifecycle.OnCancellation = handler }
}

// OptJobOnBroken sets a lifecycle hook.
func OptJobOnBroken(handler func(context.Context)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobLifecycle.OnBroken = handler }
}

// OptJobOnFixed sets a lifecycle hook.
func OptJobOnFixed(handler func(context.Context)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobLifecycle.OnFixed = handler }
}

// OptJobOnEnabled sets a lifecycle hook.
func OptJobOnEnabled(handler func(context.Context)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobLifecycle.OnEnabled = handler }
}

// OptJobOnDisabled sets a lifecycle hook.
func OptJobOnDisabled(handler func(context.Context)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobLifecycle.OnDisabled = handler }
}

// OptJobRestoreHistory sets a lifecycle hook.
func OptJobRestoreHistory(implementation func(context.Context) ([]JobInvocation, error)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobLifecycle.RestoreHistory = implementation }
}

// OptJobPersistHistory sets a lifecycle hook.
func OptJobPersistHistory(implementation func(context.Context, []JobInvocation) error) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobLifecycle.PersistHistory = implementation }
}

// JobBuilder allows for job creation w/o a fully formed struct.
type JobBuilder struct {
	JobName             string
	JobConfig           JobConfig
	JobLifecycle        JobLifecycle
	JobAction           Action
	JobScheduleProvider func() Schedule
}

// Name returns the job name.
func (jb *JobBuilder) Name() string {
	return jb.JobName
}

// Schedule returns the job schedule if a provider is set.
func (jb *JobBuilder) Schedule() Schedule {
	if jb.JobScheduleProvider != nil {
		return jb.JobScheduleProvider()
	}
	return nil
}

// Config returns the job config.
func (jb *JobBuilder) Config() JobConfig {
	return jb.JobConfig
}

// Lifecycle returns the job lifecycle hooks.
func (jb *JobBuilder) Lifecycle() JobLifecycle {
	return jb.JobLifecycle
}

// Execute runs the job action if it's set.
func (jb *JobBuilder) Execute(ctx context.Context) error {
	if jb.JobAction != nil {
		return jb.JobAction(ctx)
	}
	return nil
}
