package cron

// NOTE: ALL TIMES ARE IN UTC. JUST USE UTC.

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
)

// New returns a new job manager.
func New(options ...JobManagerOption) *JobManager {
	jm := JobManager{
		Latch: async.NewLatch(),
		Jobs:  map[string]*JobScheduler{},
	}
	for _, option := range options {
		option(&jm)
	}
	return &jm
}

// JobManager is the main orchestration and job management object.
type JobManager struct {
	sync.Mutex
	Latch   *async.Latch
	Config  Config
	Tracer  Tracer
	Log     logger.Log
	Started time.Time
	Paused  time.Time
	Stopped time.Time
	Jobs    map[string]*JobScheduler
}

// --------------------------------------------------------------------------------
// Core Methods
// --------------------------------------------------------------------------------

// LoadJobs loads a variadic list of jobs.
func (jm *JobManager) LoadJobs(jobs ...Job) error {
	jm.Lock()
	defer jm.Unlock()

	for _, job := range jobs {
		jobName := job.Name()
		if _, hasJob := jm.Jobs[jobName]; hasJob {
			return ex.New(ErrJobAlreadyLoaded, ex.OptMessagef("job: %s", job.Name()))
		}

		scheduler := NewJobScheduler(job,
			OptJobSchedulerTracer(jm.Tracer),
			OptJobSchedulerLog(jm.Log),
			OptJobSchedulerConfig(jm.Config),
		)
		if typed, ok := job.(HistoryPersister); ok {
			var err error
			scheduler.History, err = typed.HistoryRestore(context.Background())
			if err != nil {
				logger.MaybeError(jm.Log, err)
			}
		}
		jm.Jobs[jobName] = scheduler
	}
	return nil
}

// DisableJobs disables a variadic list of job names.
func (jm *JobManager) DisableJobs(jobNames ...string) error {
	jm.Lock()
	defer jm.Unlock()

	for _, jobName := range jobNames {
		if job, ok := jm.Jobs[jobName]; ok {
			job.Disable()
		} else {
			return ex.New(ErrJobNotFound, ex.OptMessagef("job: %s", jobName))
		}
	}
	return nil
}

// EnableJobs enables a variadic list of job names.
func (jm *JobManager) EnableJobs(jobNames ...string) error {
	jm.Lock()
	defer jm.Unlock()

	for _, jobName := range jobNames {
		if job, ok := jm.Jobs[jobName]; ok {
			job.Enable()
		} else {
			return ex.New(ErrJobNotFound, ex.OptMessagef("job: %s", jobName))
		}
	}
	return nil
}

// HasJob returns if a jobName is loaded or not.
func (jm *JobManager) HasJob(jobName string) (hasJob bool) {
	jm.Lock()
	defer jm.Unlock()
	_, hasJob = jm.Jobs[jobName]
	return
}

// Job returns a job metadata by name.
func (jm *JobManager) Job(jobName string) (job *JobScheduler, err error) {
	jm.Lock()
	defer jm.Unlock()
	if jobScheduler, hasJob := jm.Jobs[jobName]; hasJob {
		job = jobScheduler
	} else {
		err = ex.New(ErrJobNotLoaded, ex.OptMessagef("job: %s", jobName))
	}
	return
}

// IsJobDisabled returns if a job is disabled.
func (jm *JobManager) IsJobDisabled(jobName string) (value bool) {
	jm.Lock()
	defer jm.Unlock()

	if job, hasJob := jm.Jobs[jobName]; hasJob {
		value = job.Disabled
		if job.EnabledProvider != nil {
			value = value || !job.EnabledProvider()
		}
	}
	return
}

// IsJobRunning returns if a job is currently running.
func (jm *JobManager) IsJobRunning(jobName string) (isRunning bool) {
	jm.Lock()
	defer jm.Unlock()

	if job, ok := jm.Jobs[jobName]; ok {
		isRunning = len(job.Current) > 0
	}
	return
}

// IsJobInvocationRunning returns if a job invocation is currently running.
func (jm *JobManager) IsJobInvocationRunning(jobName, invocationID string) (isRunning bool) {
	jm.Lock()
	defer jm.Unlock()

	if job, ok := jm.Jobs[jobName]; ok {
		_, isRunning = job.Current[invocationID]
	}
	return
}

// RunJobs runs a variadic list of job names.
func (jm *JobManager) RunJobs(jobNames ...string) error {
	jm.Lock()
	defer jm.Unlock()

	for _, jobName := range jobNames {
		if job, ok := jm.Jobs[jobName]; ok {
			job.Run()
		} else {
			return ex.New(ErrJobNotLoaded, ex.OptMessagef("job: %s", jobName))
		}
	}
	return nil
}

// RunJob runs a job by jobName on demand.
func (jm *JobManager) RunJob(jobName string) error {
	jm.Lock()
	defer jm.Unlock()

	job, ok := jm.Jobs[jobName]
	if !ok {
		return ex.New(ErrJobNotLoaded, ex.OptMessagef("job: %s", jobName))
	}
	go job.Run()
	return nil
}

// RunAllJobs runs every job that has been loaded in the JobManager at once.
func (jm *JobManager) RunAllJobs() {
	jm.Lock()
	defer jm.Unlock()

	for _, job := range jm.Jobs {
		go job.Run()
	}
}

// CancelJob cancels (sends the cancellation signal) to a running job.
func (jm *JobManager) CancelJob(jobName string) (err error) {
	jm.Lock()
	defer jm.Unlock()

	jobScheduler, ok := jm.Jobs[jobName]
	if !ok {
		err = ex.New(ErrJobNotFound, ex.OptMessagef("job: %s", jobName))
		return
	}
	jobScheduler.Cancel()
	return
}

// Status returns a status object.
func (jm *JobManager) Status() *Status {
	jm.Lock()
	defer jm.Unlock()

	var jobManagerStatus JobManagerStatus
	if jm.Latch.IsStarted() {
		jobManagerStatus = JobManagerStatusRunning
	} else if jm.Latch.IsPaused() {
		jobManagerStatus = JobManagerStatusPaused
	} else if jm.Latch.IsResuming() {
		jobManagerStatus = JobManagerStatusResuming
	} else {
		jobManagerStatus = JobManagerStatusStopped
	}

	status := Status{
		Status:  jobManagerStatus,
		Started: jm.Started,
		Stopped: jm.Stopped,
		Running: map[string][]*JobInvocation{},
	}

	for _, job := range jm.Jobs {
		status.Jobs = append(status.Jobs, job)
		if job.Last != nil {
			if job.Last.Started.After(status.JobLastStarted) {
				status.JobLastStarted = job.Last.Started
			}
		}
		for _, ji := range job.Current {
			status.Running[job.Name] = append(status.Running[job.Name], ji)
		}
	}
	sort.Sort(JobSchedulersByJobNameAsc(status.Jobs))
	return &status
}

//
// Life Cycle
//

// Start starts the job manager and blocks.
func (jm *JobManager) Start() error {
	if err := jm.StartAsync(); err != nil {
		return err
	}
	<-jm.Latch.NotifyStopped()
	return nil
}

// NotifyStarted implements graceful.Graceful.
func (jm *JobManager) NotifyStarted() <-chan struct{} {
	return jm.Latch.NotifyStarted()
}

// StartAsync starts the job manager and the loaded jobs.
// It does not block.
func (jm *JobManager) StartAsync() error {
	if !jm.Latch.CanStart() {
		return fmt.Errorf("already started")
	}
	jm.Latch.Starting()
	logger.MaybeInfo(jm.Log, "job manager starting")
	for _, job := range jm.Jobs {
		job.Log = jm.Log
		job.Tracer = jm.Tracer
		job.Config = jm.Config
		go job.Start()
		<-job.NotifyStarted()
	}
	jm.Latch.Started()
	jm.Started = time.Now().UTC()
	logger.MaybeInfo(jm.Log, "job manager started")
	return nil
}

// Pause stops the schedule runner for a JobManager.
func (jm *JobManager) Pause() error {
	if !jm.Latch.CanPause() {
		return fmt.Errorf("already paused")
	}
	jm.Latch.Pausing()
	logger.MaybeInfo(jm.Log, "job manager pausing")
	for _, job := range jm.Jobs {
		job.Stop()
	}
	jm.Latch.Paused()
	jm.Paused = time.Now().UTC()
	logger.MaybeInfo(jm.Log, "job manager pausing complete")
	return nil
}

// Resume stops the schedule runner for a JobManager.
func (jm *JobManager) Resume() error {
	if !jm.Latch.CanStart() {
		return fmt.Errorf("already resumed")
	}
	jm.Latch.Resuming()
	logger.MaybeInfo(jm.Log, "job manager resuming")
	for _, job := range jm.Jobs {
		go job.Start()
		<-job.NotifyStarted()
	}
	jm.Latch.Paused()
	jm.Started = time.Now().UTC()
	logger.MaybeInfo(jm.Log, "job manager resuming complete")
	return nil
}

// Stop stops the schedule runner for a JobManager.
func (jm *JobManager) Stop() error {
	if !jm.Latch.CanStop() {
		return fmt.Errorf("already stopped")
	}
	jm.Latch.Stopping()
	logger.MaybeInfo(jm.Log, "job manager shutting down")
	for _, job := range jm.Jobs {
		job.Stop()
	}
	jm.Latch.Stopped()
	jm.Stopped = time.Now().UTC()
	logger.MaybeInfo(jm.Log, "job manager shutdown complete")
	return nil
}

// NotifyStopped implements graceful.Graceful.
func (jm *JobManager) NotifyStopped() <-chan struct{} {
	return jm.Latch.NotifyStopped()
}
