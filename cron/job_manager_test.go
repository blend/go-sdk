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
	logger "github.com/blend/go-sdk/logger"
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
	a := assert.New(t)
	a.StartTimeout(2000 * time.Millisecond)
	defer a.EndTimeout()

	manager := New()
	job := &testWithDisabled{
		disabled: false,
		action:   func() {},
	}

	manager.LoadJobs(job)
	a.False(manager.IsJobDisabled("testWithEnabled"))
	manager.DisableJobs("testWithEnabled")
	a.True(manager.IsJobDisabled("testWithEnabled"))
	job.disabled = true
	a.True(manager.IsJobDisabled("testWithEnabled"))
	manager.EnableJobs("testWithEnabled")
	a.True(manager.IsJobDisabled("testWithEnabled"))
}

func TestFiresErrorOnTaskError(t *testing.T) {
	a := assert.New(t)
	a.StartTimeout(2000 * time.Millisecond)
	defer a.EndTimeout()

	agent := logger.All(logger.OptOutput(ioutil.Discard))
	defer agent.Close()

	manager := New(OptLog(agent))
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

	manager.LoadJobs(NewJob(OptJobName("error_test"), OptJobAction(func(ctx context.Context) error {
		defer wg.Done()
		return fmt.Errorf("this is only a test")
	})))
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
		OnFinish: func(ctx context.Context) {
			defer wg.Done()
			didCallFinish = true
			finishTaskCorrect = GetJobInvocation(ctx).JobName == "tracer-test"
			errorUnset = GetJobInvocation(ctx).Err == nil
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

func TestJobManagerRunJobs(t *testing.T) {
	assert := assert.New(t)

	jm := New()
	assert.Nil(jm.StartAsync())
	defer jm.Stop()

	job0 := uuid.V4().String()
	job1 := uuid.V4().String()
	job2 := uuid.V4().String()

	var job0ran bool
	var job1ran bool
	var job2ran bool

	wg := sync.WaitGroup{}
	wg.Add(2)
	assert.Nil(jm.LoadJobs(NewJob(OptJobName(job0), OptJobAction(func(_ context.Context) error {
		defer wg.Done()
		job0ran = true
		return nil
	}))))
	assert.Nil(jm.LoadJobs(NewJob(OptJobName(job1), OptJobAction(func(_ context.Context) error {
		defer wg.Done()
		job1ran = true
		return nil
	}))))
	assert.Nil(jm.LoadJobs(NewJob(OptJobName(job2), OptJobAction(func(_ context.Context) error {
		defer wg.Done()
		job2ran = true
		return nil
	}))))

	assert.Nil(jm.RunJobs(job0, job2))
	wg.Wait()

	assert.True(job0ran)
	assert.False(job1ran)
	assert.True(job2ran)
}

func TestJobManagerRunAllJobs(t *testing.T) {
	assert := assert.New(t)

	jm := New()
	assert.Nil(jm.StartAsync())
	defer jm.Stop()

	job0 := uuid.V4().String()
	job1 := uuid.V4().String()
	job2 := uuid.V4().String()

	var job0ran bool
	var job1ran bool
	var job2ran bool

	wg := sync.WaitGroup{}
	wg.Add(3)
	assert.Nil(jm.LoadJobs(NewJob(OptJobName(job0), OptJobAction(func(_ context.Context) error {
		defer wg.Done()
		job0ran = true
		return nil
	}))))
	assert.Nil(jm.LoadJobs(NewJob(OptJobName(job1), OptJobAction(func(_ context.Context) error {
		defer wg.Done()
		job1ran = true
		return nil
	}))))
	assert.Nil(jm.LoadJobs(NewJob(OptJobName(job2), OptJobAction(func(_ context.Context) error {
		defer wg.Done()
		job2ran = true
		return nil
	}))))

	jm.RunAllJobs()
	wg.Wait()

	assert.True(job0ran)
	assert.True(job1ran)
	assert.True(job2ran)
}

func TestJobManagerJobLifecycle(t *testing.T) {
	assert := assert.New(t)

	jm := New()
	assert.Nil(jm.StartAsync())
	defer jm.Stop()

	var shouldFail bool
	j := newBrokenFixedTest(func(_ context.Context) error {
		defer func() {
			shouldFail = !shouldFail
		}()
		if shouldFail {
			return fmt.Errorf("only a test")
		}
		return nil
	})
	jm.LoadJobs(j)

	assert.Nil(jm.RunJob("broken-fixed"))
	assert.Nil(jm.RunJob("broken-fixed"))
	<-j.BrokenSignal
	assert.Nil(jm.RunJob("broken-fixed"))
	<-j.FixedSignal

	assert.Equal(3, j.Starts)
	assert.Equal(1, j.Failures)
	assert.Equal(2, j.Completes)
}

func TestJobManagerJob(t *testing.T) {
	assert := assert.New(t)

	jm := New()
	j := newBrokenFixedTest(func(_ context.Context) error {
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

	meta, err := jm.Job("load-job-test-minimum")
	assert.Nil(err)
	assert.NotNil(meta)

	assert.Equal("load-job-test-minimum", meta.Name())
	assert.NotNil(meta.Job)

	assert.NotNil(meta.DisabledProvider)
	assert.Equal(DefaultDisabled, meta.DisabledProvider())
	assert.NotNil(meta.TimeoutProvider)
	assert.Zero(meta.TimeoutProvider())
	assert.NotNil(meta.ShouldSkipLoggerListenersProvider)
	assert.Equal(DefaultShouldSkipLoggerListeners, meta.ShouldSkipLoggerListenersProvider())
	assert.NotNil(meta.ShouldSkipLoggerOutputProvider)
	assert.Equal(DefaultShouldSkipLoggerOutput, meta.ShouldSkipLoggerOutputProvider())

	jm.LoadJobs(&testJobWithTimeout{TimeoutDuration: time.Second})

	meta, err = jm.Job("testJobWithTimeout")
	assert.Nil(err)
	assert.NotNil(meta)
	assert.Equal(time.Second, meta.TimeoutProvider())
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

func TestJobManagerStatus(t *testing.T) {
	assert := assert.New(t)

	jm := New()

	jm.LoadJobs(NewJob(OptJobName("test-0")))
	jm.LoadJobs(NewJob(OptJobName("test-1")))

	checked := make(chan struct{})
	proceed := make(chan struct{})
	jm.LoadJobs(NewJob(OptJobName("is-running-test"), OptJobAction(func(_ context.Context) error {
		close(proceed)
		<-checked
		return nil
	})))

	jm.RunJob("is-running-test")
	<-proceed
	status := jm.Status()
	close(checked)

	assert.NotNil(status)
	assert.Len(status.Jobs, 3)
	assert.Len(status.Running, 1)
}

func TestJobManagerCancelJob(t *testing.T) {
	assert := assert.New(t)

	jm := New()

	proceed := make(chan struct{})
	cancelled := make(chan struct{})
	jm.LoadJobs(NewJob(OptJobName("is-running-test"), OptJobAction(func(ctx context.Context) error {
		close(proceed)
		select {
		case <-ctx.Done():
			return nil
		}
	}), OptJobOnCancellation(func(_ *JobInvocation) {
		close(cancelled)
	})))

	jm.RunJob("is-running-test")
	<-proceed
	assert.Nil(jm.CancelJob("is-running-test"))
	<-cancelled
	assert.False(jm.IsJobRunning("is-running-test"))
}

func TestJobManagerStatusRunning(t *testing.T) {
	assert := assert.New(t)

	jobDidRun := make(chan struct{})
	jobStarted := make(chan struct{})
	jobShouldProceed := make(chan struct{})
	jm := New()
	jm.LoadJobs(NewJob(OptJobName("status-running-test"), OptJobAction(func(_ context.Context) error {
		defer close(jobDidRun)
		close(jobStarted)
		<-jobShouldProceed
		return nil
	})))

	status := jm.Status()
	assert.Empty(status.Running)
	jm.RunJob("status-running-test")
	<-jobStarted
	status = jm.Status()
	assert.Len(status.Running, 1)
	close(jobShouldProceed)
	<-jobDidRun
}

func TestJobManagerEnableDisableJob(t *testing.T) {
	assert := assert.New(t)

	name := "enable-disable-test"
	jm := New()
	jm.LoadJobs(NewJob(OptJobName(name)))

	j, err := jm.Job(name)
	assert.Nil(err)
	assert.False(j.Disabled)

	jm.DisableJobs(name)
	j, err = jm.Job(name)
	assert.Nil(err)
	assert.True(j.Disabled)

	jm.EnableJobs(name)
	j, err = jm.Job(name)
	assert.Nil(err)
	assert.False(j.Disabled)
}
