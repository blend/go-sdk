package cron

import (
	"context"
	"fmt"
	"sync"
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

	js := NewJobScheduler(NewJob())
	js.JobConfig.HistoryMaxCount = ref.Int(10)
	js.JobConfig.HistoryMaxAge = ref.Duration(6 * time.Hour)

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
		OptJobHistoryEnabled(func() bool { return true }),
		OptJobHistoryPersistenceEnabled(func() bool { return true }),
		OptJobHistoryMaxCount(func() int { return 5 }),
		OptJobHistoryMaxAge(func() time.Duration { return 6 * time.Hour }),
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
		OptJobHistoryEnabled(ConstBool(true)),
		OptJobHistoryPersistenceEnabled(ConstBool(true)),
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

	job.HistoryEnabledProvider = ConstBool(false)

	js.Run()
	assert.Len(history, 3)
	assert.Len(js.History, 5)

	assert.Nil(js.RestoreHistory(context.Background()))
	assert.Len(js.History, 3)

	job.HistoryEnabledProvider = ConstBool(true)

	job.PersistHistoryHandler = func(_ context.Context, h []JobInvocation) error {
		return fmt.Errorf("only a test")
	}
	assert.NotNil(js.PersistHistory(context.Background()))
}

func TestJobSchedulerLabels(t *testing.T) {
	assert := assert.New(t)

	job := NewJob(OptJobName("test"), OptJobAction(noop))
	js := NewJobScheduler(job)
	js.Last = &JobInvocation{
		State: JobInvocationStateComplete,
	}
	labels := js.Labels()
	assert.Equal("test", labels["name"])

	job.LabelsProvider = func() map[string]string {
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

func TestSchedulerImediatelyThen(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(6)
	var runCount int
	var durations []time.Duration
	last := time.Now().UTC()
	job := NewJob(
		OptJobSchedule(Immediately().Then(Times(5, Every(time.Millisecond)))),
		OptJobAction(func(ctx context.Context) error {
			defer wg.Done()
			now := time.Now().UTC()
			runCount++
			durations = append(durations, now.Sub(last))
			last = now
			return nil
		}),
	)
	js := NewJobScheduler(job)
	go js.Start()
	<-js.NotifyStarted()

	wg.Wait()
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
