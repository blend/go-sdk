/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync/atomic"
	"time"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/logger"
)

const (
	// N is the number of jobs to load.
	N = 32

	// Q is the total simulation time.
	Q = 10 * time.Second

	// JobRunEvery is the job interval.
	JobRunEvery = 5 * time.Second

	// JobTimeout is the timeout for the jobs.
	JobTimeout = 3 * time.Second

	// JobShortRunTime is the short run time.
	JobShortRunTime = 2 * time.Second

	// JobLongRunTime is the long run time (will induce a timeout.)
	JobLongRunTime = 8 * time.Second
)

var startedCount int32
var completeCount int32
var expectedTimeoutCount int32
var timeoutCount int32

var (
	_ cron.Job               = (*loadTestJob)(nil)
	_ cron.ScheduleProvider  = (*loadTestJob)(nil)
	_ cron.ConfigProvider    = (*loadTestJob)(nil)
	_ cron.LifecycleProvider = (*loadTestJob)(nil)
)

type loadTestJob struct {
	id      int
	running bool
}

func (j *loadTestJob) Name() string {
	return fmt.Sprintf("loadTestJob_%d", j.id)
}

// Config returns a job config.
func (j *loadTestJob) Config() cron.JobConfig {
	return cron.JobConfig{
		Timeout: JobTimeout,
	}
}

// Lifecycle implements cron.LifecycleProvider.
func (j *loadTestJob) Lifecycle() cron.JobLifecycle {
	return cron.JobLifecycle{
		OnCancellation: j.OnCancellation,
	}
}

func (j *loadTestJob) Execute(ctx context.Context) error {
	atomic.AddInt32(&startedCount, 1)
	j.running = true

	var runFor time.Duration
	var randValue = rand.Float64()
	if randValue <= 0.5 { // 50% split between short vs. long.
		runFor = JobShortRunTime
	} else {
		atomic.AddInt32(&expectedTimeoutCount, 1)
		runFor = JobLongRunTime
	}

	runForElapsed := time.After(runFor)
	select {
	case <-runForElapsed:
		j.running = false
		atomic.AddInt32(&completeCount, 1)
		return nil
	case <-ctx.Done():
		j.running = false
		return nil
	}
}

func (j *loadTestJob) OnCancellation(_ context.Context) {
	atomic.AddInt32(&timeoutCount, 1)
	j.running = false
}

func (j *loadTestJob) Status() string {
	if j.running {
		return "Request in progress."
	}
	return "Request idle."
}

func (j *loadTestJob) Schedule() cron.Schedule {
	return cron.Every(JobRunEvery)
}

func main() {
	jm := cron.New(
		cron.OptLog(logger.Prod()),
	)
	defer func() {
		jm.Stop()
	}()

	if JobLongRunTime < JobTimeout {
		fmt.Printf("Long Run Time: %v is less than the Time Out: %v\n", JobTimeout, JobLongRunTime)
		fmt.Printf("This will cause the Completed vs. Timed Out counts to be wrong.\n")
		os.Exit(1)
	}

	for x := 0; x < N; x++ {
		jm.LoadJobs(&loadTestJob{id: x})
	}
	fmt.Printf("Loaded %d Job Instances.\n", N)
	fmt.Printf("Jobs run every %v\n", JobRunEvery)
	fmt.Printf("Jobs run for %v/%v\n", JobShortRunTime, JobLongRunTime)
	fmt.Printf("Jobs timeout %v\n", JobTimeout)
	fmt.Println()

	if err := jm.StartAsync(); err != nil {
		logger.FatalExit(err)
	}

	time.Sleep(Q)

	if err := jm.Stop(); err != nil {
		fmt.Fprintf(os.Stderr, "error stopping job manager: %+v\n", err)
		os.Exit(1)
	}

	// given 30 seconds total
	// and running every 5 seconds
	// we expect each job to run 5 times (ish)

	expectedStarted := N * ((int64(Q) / int64(JobRunEvery)) - 1)
	expectedCompleted := expectedStarted >> 1

	fmt.Printf("Expected Jobs Started:   %d\n", expectedStarted)
	fmt.Printf("Actual Jobs Started:     %d\n\n", startedCount)

	fmt.Printf("Expected Jobs Completed: %d\n", expectedCompleted)
	fmt.Printf("Actual Jobs Completed:   %d\n\n", completeCount)

	fmt.Printf("Expected Jobs Timed Out: %d\n", expectedTimeoutCount)
	fmt.Printf("Actual Jobs Timed Out:   %d\n", timeoutCount)
}
