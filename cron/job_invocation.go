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
		Started: Now(),
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
	ID      string `json:"id"`
	JobName string `json:"jobName"`

	Started   time.Time `json:"started"`
	Complete  time.Time `json:"complete"`
	Cancelled time.Time `json:"cancelled"`
	Timeout   time.Time `json:"timeout"`
	Errored   time.Time `json:"errored"`
	Err       error     `json:"-"`

	Parameters JobParameters       `json:"parameters"`
	Status     JobInvocationStatus `json:"status"`
	State      interface{}         `json:"state"`

	// these cannot be json marshaled.
	Context context.Context    `json:"-"`
	Cancel  context.CancelFunc `json:"-"`
	Done    chan struct{}      `json:"-"`
}

// Elapsed returns the elapsed time for the invocation.
func (ji JobInvocation) Elapsed() time.Duration {
	if !ji.Complete.IsZero() {
		return ji.Complete.Sub(ji.Started)
	} else if !ji.Cancelled.IsZero() {
		return ji.Cancelled.Sub(ji.Started)
	} else if !ji.Timeout.IsZero() {
		return ji.Timeout.Sub(ji.Started)
	} else if !ji.Errored.IsZero() {
		return ji.Errored.Sub(ji.Started)
	}
	return 0
}
