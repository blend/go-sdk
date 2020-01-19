package cron

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/graceful"
)

var (
	_ graceful.Graceful = (*JobScheduler)(nil)
)

func TestJobSchedulerEnableDisable(t *testing.T) {
	assert := assert.New(t)

	var triggerdOnEnabled, triggeredOnDisabled bool
	js := NewJobScheduler(
		NewJob(
			OptJobOnDisabled(func(_ context.Context) { triggeredOnDisabled = true }),
			OptJobOnEnabled(func(_ context.Context) { triggerdOnEnabled = true }),
		),
	)

	js.Disable()
	assert.True(js.Disabled())
	assert.False(js.CanBeScheduled())
	assert.True(triggeredOnDisabled)

	js.Enable()
	assert.False(js.Disabled())
	assert.True(js.CanBeScheduled())
	assert.True(triggerdOnEnabled)
}

func TestJobSchedulerLabels(t *testing.T) {
	assert := assert.New(t)

	job := NewJob(OptJobName("test"), OptJobAction(noop))
	js := NewJobScheduler(job)
	js.Last = &JobInvocation{
		Status: JobInvocationStatusSuccess,
	}
	labels := js.Labels()
	assert.Equal("test", labels["name"])

	job.JobConfig.Labels = map[string]string{
		"name": "not-test",
		"foo":  "bar",
		"fuzz": "wuzz",
	}

	labels = js.Labels()
	assert.Equal("true", labels["enabled"])
	assert.Equal("false", labels["active"])
	assert.Equal("not-test", labels["name"])
	assert.Equal("bar", labels["foo"])
	assert.Equal("wuzz", labels["fuzz"])
	assert.Equal(JobInvocationStatusSuccess, labels["last"])
}

func TestJobSchedulerJobParameters(t *testing.T) {
	assert := assert.New(t)

	var contextParameters, invocationParameters JobParameters

	done := make(chan struct{})
	js := NewJobScheduler(
		NewJob(
			OptJobName("test"),
			OptJobAction(func(ctx context.Context) error {
				defer close(done)
				ji := GetJobInvocation(ctx)
				invocationParameters = ji.Parameters
				contextParameters = GetJobParameters(ctx)
				return nil
			}),
		),
	)

	testParameters := JobParameters{
		"foo":    "bar",
		"moo":    "loo",
		"bailey": "dog",
	}

	ji, err := js.RunAsyncContext(WithJobParameters(context.Background(), testParameters))
	assert.Nil(err)
	assert.Equal(testParameters, ji.Parameters)
	<-done
	assert.Equal(testParameters, contextParameters)
	assert.Equal(testParameters, invocationParameters)
}
