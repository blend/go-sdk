package cron

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/uuid"
)

var (
	_ graceful.Graceful = (*JobScheduler)(nil)
)

func TestJobSchedulerCullHistoryMaxAge(t *testing.T) {
	assert := assert.New(t)

	js := NewJobScheduler(NewJob())
	js.HistoryMaxCountProvider = func() int { return 10 }
	js.HistoryMaxAgeProvider = func() time.Duration { return 6 * time.Hour }

	js.History = []JobInvocation{
		{ID: uuid.V4().String(), Started: time.Now().Add(-10 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-9 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-8 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-7 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-6 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-5 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-4 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-3 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-2 * time.Hour)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-1 * time.Hour)},
	}

	filtered := js.cullHistory()
	assert.Len(filtered, 5)
}

func TestJobSchedulerCullHistoryMaxCount(t *testing.T) {
	assert := assert.New(t)

	js := NewJobScheduler(NewJob())
	js.HistoryMaxCountProvider = func() int { return 5 }
	js.HistoryMaxAgeProvider = func() time.Duration { return 6 * time.Hour }

	js.History = []JobInvocation{
		{ID: uuid.V4().String(), Started: time.Now().Add(-10 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-9 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-8 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-7 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-6 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-5 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-4 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-3 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-2 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-1 * time.Minute)},
	}

	filtered := js.cullHistory()
	assert.Len(filtered, 5)
}

func TestJobSchedulerJobInvocation(t *testing.T) {
	assert := assert.New(t)

	id7 := uuid.V4().String()

	js := NewJobScheduler(NewJob())
	js.History = []JobInvocation{
		{ID: uuid.V4().String(), Started: time.Now().Add(-10 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-9 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-8 * time.Minute)},
		{ID: id7, Started: time.Now().Add(-7 * time.Minute), Err: fmt.Errorf("this is a test")},
		{ID: uuid.V4().String(), Started: time.Now().Add(-6 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-5 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-4 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-3 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-2 * time.Minute)},
		{ID: uuid.V4().String(), Started: time.Now().Add(-1 * time.Minute)},
	}

	ji := js.JobInvocation(id7)
	assert.NotNil(ji.Err)
}

func TestJobSchedulerEnableDisable(t *testing.T) {
	assert := assert.New(t)

	var triggerdOnEnabled, triggeredOnDisabled bool

	js := NewJobScheduler(
		NewJob(
			OptJobOnDisabled(func(_ context.Context) { triggeredOnDisabled = true }),
			OptJobOnEnabled(func(_ context.Context) { triggerdOnEnabled = true }),
		),
	)

	js.HistoryMaxCountProvider = func() int { return 5 }
	js.HistoryMaxAgeProvider = func() time.Duration { return 6 * time.Hour }

	js.Disable()
	assert.True(js.Disabled)
	assert.False(js.CanBeScheduled())
	assert.True(triggeredOnDisabled)

	js.Enable()
	assert.False(js.Disabled)
	assert.True(js.CanBeScheduled())
	assert.True(triggerdOnEnabled)
}

func TestJobSchedulerPersistHistory(t *testing.T) {
	assert := assert.New(t)

	js := NewJobScheduler(
		NewJob(OptJobName("foo")),
	)
	js.HistoryDisabledProvider = func() bool { return false }
	js.HistoryPersistenceEnabledProvider = func() bool { return true }

	assert.Nil(js.RestoreHistory(context.Background()))
	assert.Nil(js.PersistHistory(context.Background()))

	history := make(chan []JobInvocation, 2)
	js.HistoryPersistProvider = func(_ context.Context, h []JobInvocation) error {
		history <- h
		return nil
	}
	js.Run()
	assert.Len(<-history, 1)
	js.Run()
	assert.Len(<-history, 2)

	js.HistoryDisabledProvider = func() bool { return true }
	js.Run()
	assert.Len(<-history, 2)

	js.HistoryRestoreProvider = func(_ context.Context) ([]JobInvocation, error) {
		return []JobInvocation{
			*NewJobInvocation("foo"),
			*NewJobInvocation("foo"),
			*NewJobInvocation("foo"),
		}, nil
	}
	assert.Nil(js.RestoreHistory(context.Background()))
	assert.Len(js.History, 3)

	js.HistoryPersistProvider = func(_ context.Context, h []JobInvocation) error {
		return fmt.Errorf("only a test")
	}
	assert.NotNil(js.PersistHistory(context.Background()))
}

func TestJobSchedulerLabels(t *testing.T) {
	assert := assert.New(t)

	js := NewJobScheduler(NewJob(OptJobName("test"), OptJobAction(noop)))
	js.Last = &JobInvocation{
		State: JobInvocationStateComplete,
	}
	labels := js.Labels()
	assert.Equal("test", labels["name"])

	js.LabelsProvider = func() map[string]string {
		return map[string]string{
			"name": "not-test",
			"foo":  "bar",
			"fuzz": "wuzz",
		}
	}

	labels = js.Labels()
	assert.Equal("true", labels["enabled"])
	assert.Equal("false", labels["active"])
	assert.Equal("not-test", labels["name"])
	assert.Equal("bar", labels["foo"])
	assert.Equal("wuzz", labels["fuzz"])
	assert.Equal(JobInvocationStateComplete, labels["last"])
}

func TestJobSchedulerStats(t *testing.T) {
	assert := assert.New(t)

	js := NewJobScheduler(NewJob(OptJobName("test"), OptJobAction(noop)))
	js.History = []JobInvocation{
		{State: JobInvocationStateComplete, Elapsed: 2 * time.Second},
		{State: JobInvocationStateComplete, Elapsed: 4 * time.Second},
		{State: JobInvocationStateComplete, Elapsed: 6 * time.Second},
		{State: JobInvocationStateComplete, Elapsed: 8 * time.Second},
		{State: JobInvocationStateFailed, Elapsed: 10 * time.Second},
		{State: JobInvocationStateFailed, Elapsed: 12 * time.Second},
		{State: JobInvocationStateCancelled, Elapsed: 14 * time.Second},
		{State: JobInvocationStateCancelled, Timeout: time.Now().UTC(), Elapsed: 16 * time.Second},
	}

	stats := js.Stats()
	assert.NotNil(stats)
	assert.Equal(0.5, stats.SuccessRate)
	assert.Equal(8, stats.RunsTotal)
	assert.Equal(4, stats.RunsSuccessful)
	assert.Equal(2, stats.RunsFailed)
	assert.Equal(1, stats.RunsCancelled)
	assert.Equal(1, stats.RunsTimedOut)

	assert.Equal(16*time.Second, stats.ElapsedMax)
	assert.Equal(16*time.Second, stats.Elapsed95th)
	assert.Equal(9*time.Second, stats.Elapsed50th)
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
