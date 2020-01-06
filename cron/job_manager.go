package cron

// NOTE: ALL TIMES ARE IN UTC. JUST USE UTC.

import (
	"context"
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
		Jobs:  make(map[string]*JobScheduler),
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
		)
		if err := scheduler.RestoreHistory(context.Background()); err != nil {
			logger.MaybeError(jm.Log, err)
			continue
		}
		if typed, ok := job.(OnLoadHandler); ok {
			typed.OnLoad()
		}
		jm.Jobs[jobName] = scheduler
	}
	return nil
}

// UnloadJobs removes jobs from the manager and stops them.
func (jm *JobManager) UnloadJobs(jobNames ...string) error {
	jm.Lock()
	defer jm.Unlock()

	for _, jobName := range jobNames {
		if jobScheduler, ok := jm.Jobs[jobName]; ok {
			if typed, ok := jobScheduler.Job.(OnUnloadHandler); ok {
				typed.OnUnload()
			}
			jobScheduler.Stop()
			delete(jm.Jobs, jobName)
		} else {
			return ex.New(ErrJobNotFound, ex.OptMessagef("job: %s", jobName))
		}
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
		value = job.Disabled()
	}
	return
}

// IsJobRunning returns if a job is currently running.
func (jm *JobManager) IsJobRunning(jobName string) (isRunning bool) {
	jm.Lock()
	defer jm.Unlock()

	if job, ok := jm.Jobs[jobName]; ok {
		isRunning = !job.Idle()
	}
	return
}

// RunJob runs a job by jobName on demand.
func (jm *JobManager) RunJob(jobName string) (*JobInvocation, error) {
	jm.Lock()
	defer jm.Unlock()

	job, ok := jm.Jobs[jobName]
	if !ok {
		return nil, ex.New(ErrJobNotLoaded, ex.OptMessagef("job: %s", jobName))
	}
	return job.RunAsync()
}

// RunJobContext runs a job by jobName on demand with a given context.
func (jm *JobManager) RunJobContext(ctx context.Context, jobName string) (*JobInvocation, error) {
	jm.Lock()
	defer jm.Unlock()

	job, ok := jm.Jobs[jobName]
	if !ok {
		return nil, ex.New(ErrJobNotLoaded, ex.OptMessagef("job: %s", jobName))
	}
	return job.RunAsyncContext(ctx)
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

// State returns the job manager state.
func (jm *JobManager) State() JobManagerState {
	if jm.Latch.IsStarted() {
		if !jm.Paused.IsZero() {
			return JobManagerStatePaused
		}
		return JobManagerStateRunning
	} else if jm.Latch.IsStopped() {
		return JobManagerStateStopped
	}
	return JobManagerStateUnknown
}

// Status returns a status object.
func (jm *JobManager) Status() JobManagerStatus {
	jm.Lock()
	defer jm.Unlock()

	status := JobManagerStatus{
		State:   jm.State(),
		Started: jm.Started,
		Stopped: jm.Stopped,
		Running: map[string]JobInvocation{},
	}
	for _, job := range jm.Jobs {
		status.Jobs = append(status.Jobs, job.Status())
		if job.Last != nil {
			if job.Last.Started.After(status.JobLastStarted) {
				status.JobLastStarted = job.Last.Started
			}
		}
		if job.Current != nil {
			status.Running[job.Name()] = *job.Current
		}
	}
	sort.Sort(JobSchedulerStatusesByJobNameAsc(status.Jobs))
	return status
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

// StartAsync starts the job manager and the loaded jobs.
// It does not block.
func (jm *JobManager) StartAsync() error {
	if !jm.Latch.CanStart() {
		return async.ErrCannotStart
	}
	jm.Latch.Starting()
	logger.MaybeInfo(jm.Log, "job manager starting")
	for _, job := range jm.Jobs {
		go job.Start()
		<-job.NotifyStarted()
	}
	jm.Latch.Started()
	jm.Started = time.Now().UTC()
	logger.MaybeInfo(jm.Log, "job manager started")
	return nil
}

// Stop stops the schedule runner for a JobManager.
func (jm *JobManager) Stop() error {
	if !jm.Latch.CanStop() {
		return async.ErrCannotStop
	}
	jm.Latch.Stopping()
	logger.MaybeInfo(jm.Log, "job manager shutting down")
	for _, jobScheduler := range jm.Jobs {
		if typed, ok := jobScheduler.Job.(OnUnloadHandler); ok {
			if err := typed.OnUnload(); err != nil {
				return err
			}
		}
		jobScheduler.Stop()
	}
	jm.Latch.Stopped()
	jm.Latch.Reset()
	jm.Stopped = time.Now().UTC()
	logger.MaybeInfo(jm.Log, "job manager shutdown complete")
	return nil
}
