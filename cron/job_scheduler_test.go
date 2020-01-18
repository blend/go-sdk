package cron

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/ref"
	"github.com/blend/go-sdk/uuid"
)

var (
	_ graceful.Graceful = (*JobScheduler)(nil)
)

func TestJobSchedulerCullHistoryMaxAge(t *testing.T) {
	assert := assert.New(t)

	js := NewJobScheduler(NewJob(
		OptJobHistoryMaxCount(10),
		OptJobHistoryMaxAge(6*time.Hour),
	))
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

	js := NewJobScheduler(NewJob(
		OptJobHistoryEnabled(true),
		OptJobHistoryPersistenceEnabled(true),
		OptJobHistoryMaxCount(5),
		OptJobHistoryMaxAge(6*time.Hour),
	))

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

	ji := js.GetJobInvocationByID(id7)
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

	js.Disable()
	assert.True(js.Disabled())
	assert.False(js.CanBeScheduled())
	assert.True(triggeredOnDisabled)

	js.Enable()
	assert.False(js.Disabled())
	assert.True(js.CanBeScheduled())
	assert.True(triggerdOnEnabled)
}

func TestJobSchedulerPersistHistory(t *testing.T) {
	assert := assert.New(t)

	var history [][]JobInvocation
	job := NewJob(
		OptJobName("foo"),
		OptJobHistoryEnabled(true),
		OptJobHistoryPersistenceEnabled(true),
		OptJobPersistHistory(func(_ context.Context, h []JobInvocation) error {
			history = append(history, h)
			return nil
		}),
		OptJobRestoreHistory(func(_ context.Context) ([]JobInvocation, error) {
			return []JobInvocation{
				*NewJobInvocation("foo"),
				*NewJobInvocation("foo"),
				*NewJobInvocation("foo"),
			}, nil
		}),
	)
	assert.Empty(history)

	js := NewJobScheduler(job)

	assert.Nil(js.RestoreHistory(context.Background()))
	assert.Len(js.History, 3)
	assert.Nil(js.PersistHistory(context.Background()))

	assert.Len(history, 1)
	assert.Len(history[0], 3)

	js.Run()
	assert.Len(history[1], 4)
	assert.Len(history, 2)
	assert.Len(js.History, 4)
	js.Run()
	assert.Len(js.History, 5)
	assert.Len(history, 3)
	assert.Len(history[2], 5)

	job.JobConfig.HistoryEnabled = ref.Bool(false)

	js.Run()
	assert.Len(history, 3)
	assert.Len(js.History, 5)

	assert.Nil(js.RestoreHistory(context.Background()))
	assert.Len(js.History, 3)

	job.JobConfig.HistoryEnabled = ref.Bool(true)

	job.JobLifecycle.PersistHistory = func(_ context.Context, h []JobInvocation) error {
		return fmt.Errorf("only a test")
	}
	assert.NotNil(js.PersistHistory(context.Background()))
}

func TestJobSchedulerLabels(t *testing.T) {
	assert := assert.New(t)

	job := NewJob(OptJobName("test"), OptJobAction(noop))
	js := NewJobScheduler(job)
	js.Last = &JobInvocation{
		Status: JobInvocationStatusComplete,
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
	assert.Equal(JobInvocationStatusComplete, labels["last"])
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
