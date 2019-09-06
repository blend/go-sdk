package cron

// NOTE: ALL TIMES ARE IN UTC. JUST USE UTC.

import (
	"context"
	"fmt"
	"sort"
	"sync"

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
	*async.Latch

	Config Config
	Tracer Tracer
	Log    logger.Log
	Jobs   map[string]*JobScheduler
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

// IsJobRunning returns if a task is currently running.
func (jm *JobManager) IsJobRunning(jobName string) (isRunning bool) {
	jm.Lock()
	defer jm.Unlock()

	if job, ok := jm.Jobs[jobName]; ok {
		isRunning = job.Current != nil
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

	status := Status{
		Running: map[string][]*JobInvocation{},
	}

	for _, job := range jm.Jobs {
		status.Jobs = append(status.Jobs, job)

		if job.Current != nil {
			status.Running[job.Name] = append(status.Running[job.Name], job.Current)
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
	<-jm.NotifyStopped()
	return nil
}

// StartAsync starts the job manager and the loaded jobs.
// It does not block.
func (jm *JobManager) StartAsync() error {
	if !jm.CanStart() {
		return fmt.Errorf("already started")
	}
	jm.Starting()
	logger.MaybeInfo(jm.Log, "job manager starting")
	for _, job := range jm.Jobs {
		job.Log = jm.Log
		job.Tracer = jm.Tracer
		job.Config = jm.Config
		go job.Start()
		<-job.NotifyStarted()
	}
	jm.Started()
	logger.MaybeInfo(jm.Log, "job manager started")
	return nil
}

// Stop stops the schedule runner for a JobManager.
func (jm *JobManager) Stop() error {
	if !jm.CanStop() {
		return fmt.Errorf("already stopped")
	}
	jm.Stopping()
	logger.MaybeInfo(jm.Log, "job manager shutting down")
	for _, job := range jm.Jobs {
		job.Stop()
	}
	jm.Stopped()
	logger.MaybeInfo(jm.Log, "job manager shutdown complete")
	return nil
}
