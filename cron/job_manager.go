package cron

// NOTE: ALL TIMES ARE IN UTC. JUST USE UTC.

import (
	"fmt"
	"sort"
	"sync"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/exception"
	"github.com/blend/go-sdk/logger"
)

// New returns a new job manager.
func New() *JobManager {
	jm := JobManager{
		latch: &async.Latch{},
		jobs:  map[string]*JobScheduler{},
	}
	return &jm
}

// NewFromConfig returns a new job manager from a given config.
func NewFromConfig(cfg *Config) *JobManager {
	return New().WithConfig(cfg)
}

// NewFromEnv returns a new job manager from the environment.
func NewFromEnv() (*JobManager, error) {
	cfg, err := NewConfigFromEnv()
	if err != nil {
		return nil, err
	}
	return NewFromConfig(cfg), nil
}

// MustNewFromEnv returns a new job manager from the environment.
func MustNewFromEnv() *JobManager {
	cfg, err := NewConfigFromEnv()
	if err != nil {
		panic(err)
	}
	return NewFromConfig(cfg)
}

// JobManager is the main orchestration and job management object.
type JobManager struct {
	sync.Mutex
	latch  *async.Latch
	cfg    *Config
	tracer Tracer
	log    logger.Log
	jobs   map[string]*JobScheduler
}

// WithLogger sets the logger and returns a reference to the job manager.
func (jm *JobManager) WithLogger(log logger.Log) *JobManager {
	jm.log = log
	return jm
}

// Logger returns the diagnostics agent.
func (jm *JobManager) Logger() logger.Log {
	return jm.log
}

// Config returns the job manager config.
func (jm *JobManager) Config() *Config {
	return jm.cfg
}

// WithConfig sets the job manager config.
func (jm *JobManager) WithConfig(cfg *Config) *JobManager {
	jm.cfg = cfg
	return jm
}

// WithTracer sets the manager's tracer.
func (jm *JobManager) WithTracer(tracer Tracer) *JobManager {
	jm.tracer = tracer
	return jm
}

// Tracer returns the manager's tracer.
func (jm *JobManager) Tracer() Tracer {
	return jm.tracer
}

// Latch returns the internal latch.
func (jm *JobManager) Latch() *async.Latch {
	return jm.latch
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
		if _, hasJob := jm.jobs[jobName]; hasJob {
			return exception.New(ErrJobAlreadyLoaded).WithMessagef("job: %s", job.Name())
		}
		jm.jobs[jobName] = NewJobScheduler(jm.cfg, job).WithTracer(jm.tracer).WithLogger(jm.log)
	}
	return nil
}

// LoadJob loads a job.
func (jm *JobManager) LoadJob(job Job) error {
	jm.Lock()
	defer jm.Unlock()

	jobName := job.Name()
	if _, hasJob := jm.jobs[jobName]; hasJob {
		return exception.New(ErrJobAlreadyLoaded).WithMessagef("job: %s", job.Name())
	}
	jm.jobs[jobName] = NewJobScheduler(jm.cfg, job).WithTracer(jm.tracer).WithLogger(jm.log)
	return nil
}

// DisableJobs disables a variadic list of job names.
func (jm *JobManager) DisableJobs(jobNames ...string) error {
	jm.Lock()
	defer jm.Unlock()

	for _, jobName := range jobNames {
		if job, ok := jm.jobs[jobName]; ok {
			job.Disable()
		} else {
			return exception.New(ErrJobNotFound).WithMessagef("job: %s", jobName)
		}
	}
	return nil
}

// DisableJob stops a job from running but does not unload it.
func (jm *JobManager) DisableJob(jobName string) error {
	jm.Lock()
	defer jm.Unlock()

	job, ok := jm.jobs[jobName]
	if !ok {
		return exception.New(ErrJobNotFound).WithMessagef("job: %s", jobName)
	}
	job.Disable()
	return nil
}

// EnableJobs enables a variadic list of job names.
func (jm *JobManager) EnableJobs(jobNames ...string) error {
	jm.Lock()
	defer jm.Unlock()

	for _, jobName := range jobNames {
		if job, ok := jm.jobs[jobName]; ok {
			job.Enable()
		} else {
			return exception.New(ErrJobNotFound).WithMessagef("job: %s", jobName)
		}
	}
	return nil
}

// EnableJob enables a job that has been disabled.
func (jm *JobManager) EnableJob(jobName string) error {
	jm.Lock()
	defer jm.Unlock()
	job, ok := jm.jobs[jobName]
	if !ok {
		return exception.New(ErrJobNotFound).WithMessagef("job: %s", jobName)
	}
	job.Enable()
	return nil
}

// HasJob returns if a jobName is loaded or not.
func (jm *JobManager) HasJob(jobName string) (hasJob bool) {
	jm.Lock()
	defer jm.Unlock()
	_, hasJob = jm.jobs[jobName]
	return
}

// Job returns a job metadata by name.
func (jm *JobManager) Job(jobName string) (job *JobScheduler, err error) {
	jm.Lock()
	defer jm.Unlock()
	if jobScheduler, hasJob := jm.jobs[jobName]; hasJob {
		job = jobScheduler
	} else {
		err = exception.New(ErrJobNotLoaded).WithMessagef("job: %s", jobName)
	}
	return
}

// IsJobDisabled returns if a job is disabled.
func (jm *JobManager) IsJobDisabled(jobName string) (value bool) {
	jm.Lock()
	defer jm.Unlock()

	if job, hasJob := jm.jobs[jobName]; hasJob {
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

	if job, ok := jm.jobs[jobName]; ok {
		isRunning = job.Current != nil
	}
	return
}

// RunJobs runs a variadic list of job names.
func (jm *JobManager) RunJobs(jobNames ...string) error {
	jm.Lock()
	defer jm.Unlock()

	for _, jobName := range jobNames {
		if job, ok := jm.jobs[jobName]; ok {
			job.Run()
		} else {
			return exception.New(ErrJobNotLoaded).WithMessagef("job: %s", jobName)
		}
	}
	return nil
}

// RunJob runs a job by jobName on demand.
func (jm *JobManager) RunJob(jobName string) error {
	jm.Lock()
	defer jm.Unlock()
	job, ok := jm.jobs[jobName]
	if !ok {
		return exception.New(ErrJobNotLoaded).WithMessagef("job: %s", jobName)
	}
	go job.Run()
	return nil
}

// RunAllJobs runs every job that has been loaded in the JobManager at once.
func (jm *JobManager) RunAllJobs() {
	jm.Lock()
	defer jm.Unlock()

	for _, job := range jm.jobs {
		go job.Run()
	}
}

// CancelJob cancels (sends the cancellation signal) to a running job.
func (jm *JobManager) CancelJob(jobName string) (err error) {
	jm.Lock()
	defer jm.Unlock()

	job, ok := jm.jobs[jobName]
	if !ok {
		err = exception.New(ErrJobNotFound).WithMessagef("job: %s", jobName)
		return
	}
	job.Cancel()
	return
}

// Status returns a status object.
func (jm *JobManager) Status() *Status {
	jm.Lock()
	defer jm.Unlock()

	status := Status{
		Running: map[string][]*JobInvocation{},
	}

	for _, job := range jm.jobs {
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

// Start begins the schedule runner for a JobManager.
// It does not block.
func (jm *JobManager) Start() error {
	if !jm.latch.CanStart() {
		return fmt.Errorf("already started")
	}
	jm.latch.Starting()
	for _, job := range jm.jobs {
		job.WithTracer(jm.tracer).WithLogger(jm.log).Start()
	}
	jm.latch.Started()
	return nil
}

// Stop stops the schedule runner for a JobManager.
func (jm *JobManager) Stop() error {
	if !jm.latch.CanStop() {
		return fmt.Errorf("already stopped")
	}
	jm.latch.Stopping()
	for _, job := range jm.jobs {
		job.Stop()
	}
	jm.latch.Stopped()
	return nil
}

// NotifyStarted returns the started notification channel.
func (jm *JobManager) NotifyStarted() <-chan struct{} {
	return jm.latch.NotifyStarted()
}

// NotifyStopped returns the stopped notification channel.
func (jm *JobManager) NotifyStopped() <-chan struct{} {
	return jm.latch.NotifyStopped()
}

// IsRunning returns if the job manager is running.
// It serves as an authoritative healthcheck.
func (jm *JobManager) IsRunning() bool {
	return jm.latch.IsRunning()
}
