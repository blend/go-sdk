package cron

import (
	"context"
	"sync"
	"testing"
	"time"

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

func TestJobSchedulerJobParameterValues(t *testing.T) {
	assert := assert.New(t)

	var contextParameters JobParameters

	done := make(chan struct{})
	js := NewJobScheduler(
		NewJob(
			OptJobName("test"),
			OptJobAction(func(ctx context.Context) error {
				defer close(done)
				contextParameters = GetJobParameterValues(ctx)
				return nil
			}),
		),
	)

	testParameters := JobParameters{
		"foo":    "bar",
		"moo":    "loo",
		"bailey": "dog",
	}

	ji, err := js.RunAsyncContext(WithJobParameterValues(context.Background(), testParameters))
	assert.Nil(err)
	assert.Equal(testParameters, ji.Parameters)
	<-done
	assert.Equal(testParameters, contextParameters)
}

func TestJobSchedulerJobParameterValuesDefault(t *testing.T) {
	assert := assert.New(t)

	var contextParameters JobParameters

	defaultParameters := JobParameters{
		"bailey":  "woof",
		"default": "value",
	}

	done := make(chan struct{})
	js := NewJobScheduler(
		NewJob(
			OptJobName("test"),
			OptJobAction(func(ctx context.Context) error {
				defer close(done)
				contextParameters = GetJobInvocation(ctx).Parameters
				return nil
			}),
			OptJobConfig(JobConfig{
				ParameterValues: defaultParameters,
			}),
		),
	)
	assert.Equal("woof", js.Config().ParameterValues["bailey"])

	runParameters := JobParameters{
		"foo":    "bar",
		"moo":    "loo",
		"bailey": "dog",
	}

	ji, err := js.RunAsyncContext(WithJobParameterValues(context.Background(), runParameters))
	assert.Nil(err)
	assert.Equal("dog", ji.Parameters["bailey"])
	assert.Equal("value", ji.Parameters["default"])
	assert.Equal("bar", ji.Parameters["foo"])
	assert.Equal("loo", ji.Parameters["moo"])
	<-done
	assert.NotEmpty(contextParameters)
	assert.Equal("dog", contextParameters["bailey"])
	assert.Equal("value", contextParameters["default"])
	assert.Equal("bar", contextParameters["foo"])
	assert.Equal("loo", contextParameters["moo"])
}

func TestSchedulerImediatelyThen(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(6)
	var runCount, scheduleCallCount int
	var durations []time.Duration
	last := time.Now().UTC()
	job := scheduleProvider{
		ScheduleProvider: func() Schedule {
			scheduleCallCount++
			return Immediately().Then(Times(5, Every(time.Millisecond)))
		},
		Action: func(ctx context.Context) error {
			defer wg.Done()
			now := time.Now().UTC()
			runCount++
			durations = append(durations, now.Sub(last))
			last = now
			return nil
		},
	}
	js := NewJobScheduler(job)
	go js.Start()
	<-js.NotifyStarted()

	wg.Wait()
	assert.Equal(1, scheduleCallCount)
	assert.Equal(6, runCount)
	assert.Len(durations, 6)
	assert.True(durations[0] < time.Millisecond)
	for x := 1; x < 6; x++ {
		assert.True(durations[x] < 2*time.Millisecond, durations[x].String())
		assert.True(durations[x] > 500*time.Microsecond, durations[x].String())
	}
	assert.NotNil(js.JobSchedule)
	typed, ok := js.JobSchedule.(*ImmediateSchedule)
	assert.True(ok)
	assert.True(typed.didRun)
}
