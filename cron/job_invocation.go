package cron

import (
	"context"
	"time"
)

// NewJobInvocation returns a new job invocation.
func NewJobInvocation(jobName string) *JobInvocation {
	return &JobInvocation{
		ID:      NewJobInvocationID(),
		Started: Now(),
		Status:  JobStatusRunning,
		JobName: jobName,
	}
}

// JobInvocation is metadata for a job invocation (or instance of a job running).
type JobInvocation struct {
	ID        string             `json:"id"`
	JobName   string             `json:"jobName"`
	Started   time.Time          `json:"started"`
	Finished  time.Time          `json:"finished,omitempty"`
	Cancelled time.Time          `json:"cancelled,omitempty"`
	Timeout   time.Time          `json:"timeout,omitempty"`
	Err       error              `json:"err,omitempty"`
	Elapsed   time.Duration      `json:"elapsed"`
	Status    JobStatus          `json:"status"`
	State     interface{}        `json:"state,omitempty"`
	Context   context.Context    `json:"-"`
	Cancel    context.CancelFunc `json:"-"`
}
