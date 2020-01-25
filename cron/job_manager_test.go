package cron

import (
	"context"
	"fmt"
	"io/ioutil"
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/uuid"
)

// assert the job manager is graceful
var (
	_ graceful.Graceful = (*JobManager)(nil)
)

func TestJobManagerNew(t *testing.T) {
	assert := assert.New(t)

	jm := New()
	assert.NotNil(jm.Latch)
	assert.NotNil(jm.Jobs)
}

func TestRunJobBySchedule(t *testing.T) {
	a := assert.New(t)

	didRun := make(chan struct{})

	interval := 50 * time.Millisecond

	jm := New()
	runAt := Now().Add(interval)
	err := jm.LoadJobs(&runAtJob{
		RunAt: runAt,
		RunDelegate: func(ctx context.Context) error {
			close(didRun)
			return nil
		},
	})
	a.Nil(err)

	a.Nil(jm.StartAsync())
	defer jm.Stop()

	before := Now()
	<-didRun

	a.True(Since(before) < 2*interval)
}

func TestDisableJob(t *testing.T) {
	a := assert.New(t)

	jm := New()
	a.Nil(jm.LoadJobs(&runAtJob{RunAt: time.Now().UTC().Add(100 * time.Millisecond), RunDelegate: func(ctx context.Context) error {
		return nil
	}}))
	a.Nil(jm.DisableJobs(runAtJobName))
	a.True(jm.IsJobDisabled(runAtJobName))
}

// The goal with this test is to see if panics take down the test process or not.
func TestJobManagerJobPanicHandling(t *testing.T) {
	assert := assert.New(t)

	manager := New()
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)

	action := func(ctx context.Context) error {
		defer waitGroup.Done()
		panic("this is only a test")
	}
	manager.LoadJobs(NewJob(OptJobName("panic-test"), OptJobAction(action)))
	manager.RunJob("panic-test")
	waitGroup.Wait()
	assert.True(true, "should complete")
}

func TestEnabledProvider(t *testing.T) {
	assert := assert.New(t)

	manager := New()
	job := &testWithDisabled{
		disabled: false,
	}

	jobName := "testWithEnabled"

	manager.LoadJobs(job)

	// test provider
	assert.False(manager.IsJobDisabled(jobName))
	job.disabled = true
	assert.True(manager.IsJobDisabled(jobName))

	// test explicit
	manager.DisableJobs(jobName)
	assert.True(manager.IsJobDisabled(jobName))
	manager.EnableJobs(jobName)
	assert.False(manager.IsJobDisabled(jobName))
}

func TestFiresErrorOnTaskError(t *testing.T) {
	a := assert.New(t)

	agent := logger.All(logger.OptOutput(ioutil.Discard))
	defer agent.Close()

	manager := New(
		OptLog(agent),
	)
	defer manager.Stop()

	var errorDidFire bool
	var errorMatched bool
	wg := sync.WaitGroup{}
	wg.Add(2)

	agent.Listen(logger.Error, "foo", func(_ context.Context, e logger.Event) {
		defer wg.Done()
		errorDidFire = true
		if typed, isTyped := e.(logger.ErrorEvent); isTyped {
			if typed.Err != nil {
				errorMatched = typed.Err.Error() == "this is only a test"
			}
		}
	})
	manager.LoadJobs(NewJob(
		OptJobName("error_test"),
		OptJobAction(func(ctx context.Context) error {
			defer wg.Done()
			return fmt.Errorf("this is only a test")
		}),
	))
	manager.RunJob("error_test")
	wg.Wait()

	a.True(errorDidFire)
	a.True(errorMatched)
}

func TestManagerTracer(t *testing.T) {
	assert := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(2)
	var didCallStart, didCallFinish bool
	var startTaskCorrect, finishTaskCorrect, errorUnset bool
	manager := New(OptTracer(&mockTracer{
		OnStart: func(ctx context.Context) {
			defer wg.Done()
			didCallStart = true
			startTaskCorrect = GetJobInvocation(ctx).JobName == "tracer-test"
		},
		OnFinish: func(ctx context.Context, err error) {
			defer wg.Done()
			didCallFinish = true
			finishTaskCorrect = GetJobInvocation(ctx).JobName == "tracer-test"
			errorUnset = err == nil
		},
	}))

	manager.LoadJobs(NewJob(OptJobName("tracer-test")))
	manager.RunJob("tracer-test")
	wg.Wait()
	assert.True(didCallStart)
	assert.True(didCallFinish)
	assert.True(startTaskCorrect)
	assert.True(finishTaskCorrect)
	assert.True(errorUnset)
}

func TestJobManagerJobLifecycle(t *testing.T) {
	assert := assert.New(t)

	jm := New()
	assert.Nil(jm.StartAsync())
	defer jm.Stop()

	var shouldFail bool
	j := newLifecycleTest(func(_ context.Context) error {
		defer func() {
			shouldFail = !shouldFail
		}()
		if shouldFail {
			return fmt.Errorf("only a test")
		}
		return nil
	})
	jm.LoadJobs(j)

	completeSignal := j.CompleteSignal
	ji, err := jm.RunJob("lifecycle-test")
	assert.Nil(err)
	<-ji.Done
	<-completeSignal

	brokenSignal := j.BrokenSignal
	ji, err = jm.RunJob("lifecycle-test")
	assert.Nil(err)
	<-ji.Done
	<-brokenSignal

	fixedSignal := j.FixedSignal
	ji, err = jm.RunJob("lifecycle-test")
	assert.Nil(err)
	<-ji.Done
	<-fixedSignal

	assert.Equal(3, j.Starts)
	assert.Equal(3, j.Completes)
	assert.Equal(1, j.Failures)
	assert.Equal(2, j.Successes)
}

func TestJobManagerJob(t *testing.T) {
	assert := assert.New(t)

	jm := New()
	j := newLifecycleTest(func(_ context.Context) error {
		return nil
	})
	jm.LoadJobs(j)

	meta, err := jm.Job(j.Name())
	assert.Nil(err)
	assert.NotNil(meta)

	meta, err = jm.Job(uuid.V4().String())
	assert.NotNil(err)
	assert.Nil(meta)
}

func TestJobManagerLoadJob(t *testing.T) {
	assert := assert.New(t)

	jm := New()
	jm.LoadJobs(&loadJobTestMinimum{})

	assert.True(jm.HasJob("load-job-test-minimum"))

	jobScheduler, err := jm.Job("load-job-test-minimum")
	assert.Nil(err)
	assert.NotNil(jobScheduler)

	assert.Equal("load-job-test-minimum", jobScheduler.Name())
	assert.NotNil(jobScheduler.Job)

	assert.Equal(DefaultDisabled, jobScheduler.Disabled())
	assert.Zero(jobScheduler.Config().TimeoutOrDefault())
	assert.Equal(DefaultShouldSkipLoggerListeners, jobScheduler.Config().ShouldSkipLoggerListenersOrDefault())
	assert.Equal(DefaultShouldSkipLoggerOutput, jobScheduler.Config().ShouldSkipLoggerOutputOrDefault())

	jm.LoadJobs(&testJobWithTimeout{TimeoutDuration: time.Second})

	jobScheduler, err = jm.Job("testJobWithTimeout")
	assert.Nil(err)
	assert.NotNil(jobScheduler)
	assert.Equal(time.Second, jobScheduler.Config().TimeoutOrDefault())
}

func TestJobManagerLoadJobs(t *testing.T) {
	assert := assert.New(t)

	jm := New()
	assert.Nil(jm.LoadJobs(NewJob(OptJobName("test-0")), NewJob(OptJobName("test-1"))))
	assert.Len(jm.Jobs, 2)

	assert.NotNil(jm.LoadJobs(NewJob(OptJobName("test-0")), NewJob(OptJobName("test-1"))))
}

func TestJobManagerIsRunning(t *testing.T) {
	assert := assert.New(t)

	jm := New()

	checked := make(chan struct{})
	proceed := make(chan struct{})
	jm.LoadJobs(NewJob(OptJobName("is-running-test"), OptJobAction(func(_ context.Context) error {
		close(proceed)
		<-checked
		return nil
	})))

	jm.RunJob("is-running-test")
	<-proceed
	assert.True(jm.IsJobRunning("is-running-test"))
	close(checked)

	assert.False(jm.IsJobRunning(uuid.V4().String()))
}

func TestJobManagerCancelJob(t *testing.T) {
	assert := assert.New(t)

	jm := New()

	proceed := make(chan struct{})
	cancelled := make(chan struct{})
	jm.LoadJobs(NewJob(OptJobName("is-running-test"), OptJobAction(func(ctx context.Context) error {
		close(proceed)
		<-ctx.Done()
		return nil
	}), OptJobOnCancellation(func(_ context.Context) {
		close(cancelled)
	})))

	jm.RunJob("is-running-test")
	<-proceed
	assert.Nil(jm.CancelJob("is-running-test"))
	<-cancelled
	assert.False(jm.IsJobRunning("is-running-test"))
}

func TestJobManagerEnableDisableJob(t *testing.T) {
	assert := assert.New(t)

	name := "enable-disable-test"
	jm := New()
	jm.LoadJobs(NewJob(OptJobName(name)))

	j, err := jm.Job(name)
	assert.Nil(err)
	assert.False(j.Disabled())

	jm.DisableJobs(name)
	j, err = jm.Job(name)
	assert.Nil(err)
	assert.True(j.Disabled())

	jm.EnableJobs(name)
	j, err = jm.Job(name)
	assert.Nil(err)
	assert.False(j.Disabled())
}
