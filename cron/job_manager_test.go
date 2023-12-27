/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cron

import (
	"context"
	"fmt"
	"io"
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

func Test_JobManager_New(t *testing.T) {
	its := assert.New(t)

	jm := New(
		OptLog(logger.None()),
	)
	its.NotNil(jm.Latch)
	its.NotNil(jm.Jobs)
	its.NotNil(jm.Log)
}

func Test_JobManager_Start(t *testing.T) {
	its := assert.New(t)

	jm := New(
		OptLog(logger.None()),
	)
	defer jm.Stop()

	// start blocks, use a channel and goroutine
	errors := make(chan error, 1)
	go func() {
		if err := jm.Start(); err != nil {
			errors <- err
		}
	}()

	<-jm.Latch.NotifyStarted()
	its.Empty(errors)
	err := jm.Start()
	its.NotNil(err)
}

func Test_JobManager_DisableJobs(t *testing.T) {
	its := assert.New(t)

	jm := New()
	its.Nil(jm.LoadJobs(&runAtJob{RunAt: time.Now().UTC().Add(100 * time.Millisecond), RunDelegate: func(ctx context.Context) error {
		return nil
	}}))
	its.Nil(jm.DisableJobs(runAtJobName))
	its.True(jm.IsJobDisabled(runAtJobName))
}

func Test_JobManager_handleJobPanics(t *testing.T) {
	its := assert.New(t)

	manager := New()
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)

	action := func(ctx context.Context) error {
		defer waitGroup.Done()
		panic("this is only a test")
	}
	its.Nil(manager.LoadJobs(NewJob(OptJobName("panic-test"), OptJobAction(action))))
	_, _, err := manager.RunJob("panic-test")
	its.Nil(err)
	waitGroup.Wait()
	its.True(true, "should complete")
}

func Test_JobManager_jobConfigProvider_disabled(t *testing.T) {
	its := assert.New(t)

	manager := New()
	job := &testWithDisabled{
		disabled: false,
	}

	jobName := "testWithEnabled"

	its.Nil(manager.LoadJobs(job))

	// test provider
	its.False(manager.IsJobDisabled(jobName))
	job.disabled = true
	its.True(manager.IsJobDisabled(jobName))

	// test explicit
	its.Nil(manager.DisableJobs(jobName))
	its.True(manager.IsJobDisabled(jobName))
	its.Nil(manager.EnableJobs(jobName))
	its.False(manager.IsJobDisabled(jobName))
}

func Test_JobManager_onError(t *testing.T) {
	its := assert.New(t)

	agent := logger.All(logger.OptOutput(io.Discard))
	defer agent.Close()

	manager := New(
		OptLog(agent),
	)
	defer func() { _ = manager.Stop() }()

	var errorDidFire bool
	var errorMatched bool
	wg := sync.WaitGroup{}
	wg.Add(2)

	agent.Listen(logger.Error, uuid.V4().String(), func(_ context.Context, e logger.Event) {
		defer wg.Done()
		errorDidFire = true
		if typed, isTyped := e.(logger.ErrorEvent); isTyped {
			if typed.Err != nil {
				errorMatched = typed.Err.Error() == "this is only a test"
			}
		}
	})
	job := NewJob(
		OptJobName("error_test"),
		OptJobAction(func(ctx context.Context) error {
			defer wg.Done()
			return fmt.Errorf("this is only a test")
		}),
	)
	its.Nil(manager.LoadJobs(job))
	_, done, err := manager.RunJob(job.Name())
	its.Nil(err)
	wg.Wait()

	its.True(errorDidFire)
	its.True(errorMatched)
	<-done
}

func Test_JobManager_Tracer(t *testing.T) {
	its := assert.New(t)

	wg := sync.WaitGroup{}
	wg.Add(2)
	var didCallStart, didCallFinish bool
	var errorUnset bool
	var foundJobName string
	manager := New(OptTracer(&mockTracer{
		OnStart: func(ctx context.Context, jobName string) {
			defer wg.Done()
			didCallStart = true
			foundJobName = jobName
		},
		OnFinish: func(ctx context.Context, err error) {
			defer wg.Done()
			didCallFinish = true
			errorUnset = err == nil
		},
	}))

	its.Nil(manager.LoadJobs(NewJob(OptJobName("tracer-test"))))
	_, _, err := manager.RunJob("tracer-test")
	its.Nil(err)
	wg.Wait()
	its.True(didCallStart)
	its.True(didCallFinish)
	its.True(errorUnset)
	its.Equal("tracer-test", foundJobName)
}

func Test_JobManager_JobLifecycle(t *testing.T) {
	its := assert.New(t)

	jm := New()
	its.Nil(jm.StartAsync())
	defer func() { _ = jm.Stop() }()

	var shouldFail bool
	j := newLifecycleTest(func(ctx context.Context) error {
		defer func() {
			shouldFail = !shouldFail
		}()
		if shouldFail {
			return fmt.Errorf("only a test")
		}
		return nil
	})
	its.Nil(jm.LoadJobs(j))

	successSignal := j.SuccessSignal
	_, done, err := jm.RunJob(j.Name())
	its.Nil(err)
	<-done
	<-successSignal

	brokenSignal := j.BrokenSignal
	_, done, err = jm.RunJob(j.Name())
	its.Nil(err)
	<-done
	<-brokenSignal

	fixedSignal := j.FixedSignal
	_, done, err = jm.RunJob(j.Name())
	its.Nil(err)
	<-done
	<-fixedSignal

	its.Equal(3, j.Starts)
	its.Equal(3, j.Completes)
	its.Equal(1, j.Failures)
	its.Equal(2, j.Successes)
}

func Test_JobManager_Job(t *testing.T) {
	its := assert.New(t)

	jm := New()
	j := newLifecycleTest(func(_ context.Context) error {
		return nil
	})
	its.Nil(jm.LoadJobs(j))

	meta, err := jm.Job(j.Name())
	its.Nil(err)
	its.NotNil(meta)

	meta, err = jm.Job(uuid.V4().String())
	its.NotNil(err)
	its.Nil(meta)
}

func Test_JobManager_LoadJobs(t *testing.T) {
	its := assert.New(t)

	jm := New()
	its.Nil(jm.LoadJobs(&loadJobTestMinimum{}))

	its.True(jm.HasJob("load-job-test-minimum"))

	jobScheduler, err := jm.Job("load-job-test-minimum")
	its.Nil(err)
	its.NotNil(jobScheduler)

	its.Equal("load-job-test-minimum", jobScheduler.Name())
	its.NotNil(jobScheduler.Job)

	its.Equal(DefaultDisabled, jobScheduler.Disabled())
	its.Zero(jobScheduler.Config().TimeoutOrDefault())

	its.Nil(jm.LoadJobs(&testJobWithTimeout{TimeoutDuration: time.Second}))

	jobScheduler, err = jm.Job("testJobWithTimeout")
	its.Nil(err)
	its.NotNil(jobScheduler)
	its.Equal(time.Second, jobScheduler.Config().TimeoutOrDefault())
}

func Test_JobManager_IsRunning(t *testing.T) {
	its := assert.New(t)

	jm := New()

	checked := make(chan struct{})
	proceed := make(chan struct{})
	its.Nil(jm.LoadJobs(NewJob(OptJobName("is-running-test"), OptJobAction(func(_ context.Context) error {
		close(proceed)
		<-checked
		return nil
	}))))	// hadoooooken

	_, _, err := jm.RunJob("is-running-test")
	its.Nil(err)
	<-proceed
	its.True(jm.IsJobRunning("is-running-test"))
	close(checked)
	its.False(jm.IsJobRunning(uuid.V4().String()))
}

func Test_JobManager_CancelJob(t *testing.T) {
	its := assert.New(t)

	started := make(chan struct{})
	canceling := make(chan struct{})
	canceled := make(chan struct{})

	jm := New()
	job := NewJob(OptJobName("is-running-test"), OptJobAction(func(ctx context.Context) error {
		close(started)
		<-canceling
		time.Sleep(time.Millisecond)	// this is a pad to make the test more reliable.
		return nil
	}), OptJobOnCancellation(func(_ context.Context) {
		close(canceled)	// but signal on the lifecycle event
	}))
	its.Nil(jm.LoadJobs(job))

	_, done, err := jm.RunJob(job.Name())
	its.Nil(err)
	<-started
	close(canceling)
	its.Nil(jm.CancelJob(job.Name()))
	<-canceled
	its.False(jm.IsJobRunning(job.Name()))
	<-done
}

func Test_JobManager_EnableDisableJob(t *testing.T) {
	its := assert.New(t)

	name := "enable-disable-test"
	jm := New()
	its.Nil(jm.LoadJobs(NewJob(OptJobName(name))))

	j, err := jm.Job(name)
	its.Nil(err)
	its.False(j.Disabled())

	its.Nil(jm.DisableJobs(name))
	j, err = jm.Job(name)
	its.Nil(err)
	its.True(j.Disabled())

	its.Nil(jm.EnableJobs(name))
	j, err = jm.Job(name)
	its.Nil(err)
	its.False(j.Disabled())
}

func Test_JobManager_LoadJobs_lifecycle(t *testing.T) {
	its := assert.New(t)

	baseContext := context.WithValue(context.Background(), testContextKey{}, "load-jobs-lifecycle")
	jm := New(
		OptBaseContext(baseContext),
	)

	onLoadContexts := make(chan context.Context, 1)
	job := NewJob(
		OptJobName("load-test"),
		OptJobOnLoad(func(ctx context.Context) error {
			onLoadContexts <- ctx
			return nil
		}),
	)

	err := jm.LoadJobs(job)
	its.Nil(err)

	gotContext := <-onLoadContexts
	its.Equal("load-jobs-lifecycle", gotContext.Value(testContextKey{}))
	js := GetJobScheduler(gotContext)
	its.NotNil(js)
	its.Equal("load-test", js.Job.Name())
}

func Test_JobManager_UnloadJobs_lifecycle(t *testing.T) {
	its := assert.New(t)

	baseContext := context.WithValue(context.Background(), testContextKey{}, "load-jobs-lifecycle")
	jm := New(
		OptBaseContext(baseContext),
	)

	onUnloadContexts := make(chan context.Context, 1)
	job0 := NewJob(
		OptJobName("load-test"),
		OptJobOnUnload(func(ctx context.Context) error {
			onUnloadContexts <- ctx
			return nil
		}),
	)
	job1 := NewJob(
		OptJobName("load-test-1"),
		OptJobOnUnload(func(ctx context.Context) error {
			onUnloadContexts <- ctx
			return nil
		}),
	)

	err := jm.LoadJobs(job0, job1)
	its.Nil(err)

	its.True(jm.HasJob("load-test"))
	its.True(jm.HasJob("load-test-1"))

	err = jm.UnloadJobs("load-test")
	its.Nil(err)

	its.False(jm.HasJob("load-test"))
	its.True(jm.HasJob("load-test-1"))

	gotContext := <-onUnloadContexts
	its.Equal("load-jobs-lifecycle", gotContext.Value(testContextKey{}))
	js := GetJobScheduler(gotContext)
	its.NotNil(js)
	its.Equal("load-test", js.Job.Name())

	err = jm.UnloadJobs(uuid.V4().String())
	its.NotNil(err)
}

func Test_JobManager_Background(t *testing.T) {
	its := assert.New(t)

	jm := New()
	its.Equal(jm.Background(), context.Background())

	type contextKey struct{}
	jm = New(
		OptBaseContext(context.WithValue(context.Background(), contextKey{}, "test-value")),
	)

	ctx := jm.Background()
	its.Equal("test-value", ctx.Value(contextKey{}))
}
