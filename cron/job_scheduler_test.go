/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cron

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
)

var (
	_ graceful.Graceful = (*JobScheduler)(nil)
)

type testContextKey struct{}

func Test_JobScheduler_Background(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	js := NewJobScheduler(
		NewJob(
			OptJobName("background-test"),
			OptJobBackground(func(ctx context.Context) context.Context {
				return context.WithValue(ctx, testContextKey{}, "a test value")
			}),
		),
	)

	ctx := js.Background()
	ctx = js.withBaseContext(ctx)
	its.Equal("a test value", ctx.Value(testContextKey{}))
	foundScheduler := GetJobScheduler(ctx)
	its.NotNil(foundScheduler)
}

func Test_JobScheduler_Stop(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	js := NewJobScheduler(NewJob(
		OptJobName("stop-test"),
		OptJobSchedule(EveryHour()),
	))
	startErrors := make(chan error)
	go func() {
		startErrors <- js.Start()
	}()

	<-js.Latch.NotifyStarted()

	its.Nil(js.Stop())
	its.Nil(<-startErrors)
}

func Test_JobScheduler_EnableDisable(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	var triggerdOnEnabled, triggeredOnDisabled bool
	js := NewJobScheduler(
		NewJob(
			OptJobOnDisabled(func(_ context.Context) { triggeredOnDisabled = true }),
			OptJobOnEnabled(func(_ context.Context) { triggerdOnEnabled = true }),
		),
	)

	js.Disable()
	its.True(js.Disabled())
	its.False(js.CanBeScheduled())
	its.True(triggeredOnDisabled)

	js.Enable()
	its.False(js.Disabled())
	its.True(js.CanBeScheduled())
	its.True(triggerdOnEnabled)
}

func Test_JobScheduler_Labels(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	job := NewJob(OptJobName("test"), OptJobAction(noop))
	js := NewJobScheduler(job)
	js.last = &JobInvocation{
		Status: JobInvocationStatusSuccess,
	}
	labels := js.Labels()
	its.Equal("test", labels["name"])

	job.JobConfig.Labels = map[string]string{
		"name": "not-test",
		"foo":  "bar",
		"fuzz": "wuzz",
	}

	labels = js.Labels()
	its.Equal("true", labels["enabled"])
	its.Equal("false", labels["active"])
	its.Equal("not-test", labels["name"])
	its.Equal("bar", labels["foo"])
	its.Equal("wuzz", labels["fuzz"])
	its.Equal(JobInvocationStatusSuccess, labels["last"])
}

func Test_JobScheduler_JobParameterValues(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	var contextParameters JobParameters

	js := NewJobScheduler(
		NewJob(
			OptJobName("test"),
			OptJobAction(func(ctx context.Context) error {
				contextParameters = GetJobParameterValues(ctx)
				return nil
			}),
		),
	)

	testParameters := JobParameters{
		"foo":            "bar",
		"moo":            "loo",
		"example-string": "dog",
	}

	ji, done, err := js.RunAsyncContext(WithJobParameterValues(context.Background(), testParameters))
	its.Nil(err)
	its.Equal(testParameters, ji.Parameters)
	<-done
	its.Equal(testParameters, contextParameters)
}

func Test_JobScheduler_JobParameterValuesDefault(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	var contextParameters JobParameters

	defaultParameters := JobParameters{
		"example-string": "woof",
		"default":        "value",
	}

	js := NewJobScheduler(
		NewJob(
			OptJobName("test"),
			OptJobAction(func(ctx context.Context) error {
				contextParameters = GetJobInvocation(ctx).Parameters
				return nil
			}),
			OptJobConfig(JobConfig{
				ParameterValues: defaultParameters,
			}),
		),
	)
	its.Equal("woof", js.Config().ParameterValues["example-string"])

	runParameters := JobParameters{
		"foo":            "bar",
		"moo":            "loo",
		"example-string": "dog",
	}

	ji, done, err := js.RunAsyncContext(WithJobParameterValues(context.Background(), runParameters))
	its.Nil(err)
	its.NotNil(done)
	its.Equal("dog", ji.Parameters["example-string"])
	its.Equal("value", ji.Parameters["default"])
	its.Equal("bar", ji.Parameters["foo"])
	its.Equal("loo", ji.Parameters["moo"])
	<-done
	its.NotEmpty(contextParameters)
	its.Equal("dog", contextParameters["example-string"])
	its.Equal("value", contextParameters["default"])
	its.Equal("bar", contextParameters["foo"])
	its.Equal("loo", contextParameters["moo"])
}

func Test_JobScheduler_onJobBegin(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	var didCallLifecycleOnBegin bool
	job := NewJob(
		OptJobName("test-job"),
		OptJobOnBegin(func(_ context.Context) {
			didCallLifecycleOnBegin = true
		}),
	)
	js := NewJobScheduler(job)
	ctx := js.Background()
	ctx = js.withBaseContext(ctx)
	ctx, js.current = js.withInvocationContext(ctx)
	js.onJobBegin(ctx)

	its.False(js.current.Started.IsZero())
	its.Equal(JobInvocationStatusRunning, js.current.Status)
	its.True(didCallLifecycleOnBegin)
}

func Test_JobScheduler_onJobCompleteCanceled(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	buffer := new(bytes.Buffer)
	log := logger.Memory(
		buffer,
		logger.OptText(
			logger.OptTextHideTimestamp(),
			logger.OptTextNoColor(),
		),
	)

	var calls []string
	job := NewJob(
		OptJobName("test-job"),
		OptJobOnCancellation(func(_ context.Context) {
			calls = append(calls, "cancellation")
		}),
		OptJobOnComplete(func(_ context.Context) {
			calls = append(calls, "complete")
		}),
	)
	js := NewJobScheduler(
		job,
		OptJobSchedulerLog(log),
	)
	ctx := js.Background()
	ctx = js.withBaseContext(ctx)
	ctx, js.current = js.withInvocationContext(ctx)
	js.onJobCompleteCanceled(ctx)

	its.False(js.current.Complete.IsZero())
	its.Equal(JobInvocationStatusCanceled, js.current.Status)
	its.Equal([]string{"cancellation", "complete"}, calls)

	its.Contains(buffer.String(), "[cron.canceled]")
	its.Contains(buffer.String(), "[cron.complete]")
}

func Test_JobScheduler_onJobCompleteSuccess(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	buffer := new(bytes.Buffer)
	log := logger.Memory(
		buffer,
		logger.OptText(
			logger.OptTextHideTimestamp(),
			logger.OptTextNoColor(),
		),
	)

	var calls []string
	job := NewJob(
		OptJobName("test-job"),
		OptJobOnSuccess(func(_ context.Context) {
			calls = append(calls, "success")
		}),
		OptJobOnComplete(func(_ context.Context) {
			calls = append(calls, "complete")
		}),
	)
	js := NewJobScheduler(
		job,
		OptJobSchedulerLog(log),
	)
	ctx := js.Background()
	ctx = js.withBaseContext(ctx)
	ctx, js.current = js.withInvocationContext(ctx)
	js.onJobCompleteSuccess(ctx)

	its.False(js.current.Complete.IsZero())
	its.Equal(JobInvocationStatusSuccess, js.current.Status)
	its.Equal([]string{"success", "complete"}, calls)

	its.Contains(buffer.String(), "[cron.success]")
	its.Contains(buffer.String(), "[cron.complete]")
}

func Test_JobScheduler_onJobCompleteError(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	buffer := new(bytes.Buffer)
	log := logger.Memory(
		buffer,
		logger.OptText(
			logger.OptTextHideTimestamp(),
			logger.OptTextNoColor(),
		),
	)

	var calls []string
	job := NewJob(
		OptJobName("test-job"),
		OptJobOnError(func(_ context.Context) {
			calls = append(calls, "error")
		}),
		OptJobOnComplete(func(_ context.Context) {
			calls = append(calls, "complete")
		}),
	)
	js := NewJobScheduler(
		job,
		OptJobSchedulerLog(log),
	)
	ctx := js.Background()
	ctx = js.withBaseContext(ctx)
	ctx, js.current = js.withInvocationContext(ctx)
	js.onJobCompleteError(ctx, fmt.Errorf("this is just a test"))

	its.False(js.current.Complete.IsZero())
	its.Equal(JobInvocationStatusErrored, js.current.Status)
	its.Equal([]string{"error", "complete"}, calls)

	its.Contains(buffer.String(), "[error] this is just a test")
	its.Contains(buffer.String(), "[cron.errored]")
	its.Contains(buffer.String(), "[cron.complete]")
}
