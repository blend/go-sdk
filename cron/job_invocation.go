package cron

import (
	"context"
	"time"

	"github.com/blend/go-sdk/uuid"
)

// NewJobInvocation returns a new job invocation.
func NewJobInvocation(jobName string) *JobInvocation {
	return &JobInvocation{
		ID:      NewJobInvocationID(),
		Status:  JobInvocationStatusIdle,
		JobName: jobName,
		Done:    make(chan struct{}),
	}
}

// NewJobInvocationID returns a new pseudo-unique job invocation identifier.
func NewJobInvocationID() string {
	return uuid.V4().String()
}

// JobInvocation is metadata for a job invocation (or instance of a job running).
type JobInvocation struct {
	ID      string
	JobName string

	Started  time.Time
	Complete time.Time
	Err      error

	Parameters JobParameters
	Status     JobInvocationStatus
	State      interface{}

	// these cannot be json marshaled.
	Context context.Context
	Cancel  context.CancelFunc
	Done    chan struct{}
}

// Elapsed returns the elapsed time for the invocation.
func (ji JobInvocation) Elapsed() time.Duration {
	if !ji.Complete.IsZero() {
		return ji.Complete.Sub(ji.Started)
	}
	return 0
}
