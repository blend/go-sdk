/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package cron

import (
	"context"
	"time"

	"github.com/blend/go-sdk/ref"
	"github.com/blend/go-sdk/stringutil"
)

// Interface assertions.
var (
	_	Job			= (*JobBuilder)(nil)
	_	ScheduleProvider	= (*JobBuilder)(nil)
	_	LifecycleProvider	= (*JobBuilder)(nil)
	_	BackgroundProvider	= (*JobBuilder)(nil)
	_	ConfigProvider		= (*JobBuilder)(nil)
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

// OptJobConfig sets the job config.
func OptJobConfig(cfg JobConfig) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobConfig = cfg }
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

// OptJobOnCancellation sets the on cancellation lifecycle hook.
func OptJobOnCancellation(handler func(context.Context)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobLifecycle.OnCancellation = handler }
}

// OptJobOnSuccess sets a lifecycle hook.
func OptJobOnSuccess(handler func(context.Context)) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobLifecycle.OnSuccess = handler }
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

// OptJobOnLoad sets a lifecycle hook.
func OptJobOnLoad(handler func(context.Context) error) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobLifecycle.OnLoad = handler }
}

// OptJobOnUnload sets a lifecycle hook.
func OptJobOnUnload(handler func(context.Context) error) JobBuilderOption {
	return func(jb *JobBuilder) { jb.JobLifecycle.OnUnload = handler }
}

// OptJobBackground sets the background provider.
func OptJobBackground(provider func(context.Context) context.Context) JobBuilderOption {
	return func(jb *JobBuilder) { jb.BackgroundProvider = provider }
}

// JobBuilder allows for job creation w/o a fully formed struct.
type JobBuilder struct {
	JobName			string
	JobConfig		JobConfig
	JobLifecycle		JobLifecycle
	JobAction		Action
	JobScheduleProvider	func() Schedule
	BackgroundProvider	func(context.Context) context.Context
}

// Name returns the job name.
func (jb *JobBuilder) Name() string {
	return jb.JobName
}

// Background implements BackgroundProvider.
func (jb *JobBuilder) Background(ctx context.Context) context.Context {
	if jb.BackgroundProvider != nil {
		return jb.BackgroundProvider(ctx)
	}
	return ctx
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
