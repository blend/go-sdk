package jobkit

import (
	"context"

	"github.com/blend/go-sdk/cron"
)

// WithJobInvocationState sets the job invocation state on a context.
func WithJobInvocationState(ctx context.Context, jis *JobInvocationState) context.Context {
	ji := cron.GetJobInvocation(ctx)
	if ji == nil {
		return ctx
	}
	ji.State = jis
	return ctx
}

// GetJobInvocationState returns the job invocation state.
func GetJobInvocationState(ctx context.Context) *JobInvocationState {
	ji := cron.GetJobInvocation(ctx)
	if ji == nil {
		return nil
	}
	if typed, ok := ji.State.(*JobInvocationState); ok {
		return typed
	}
	return nil
}

// NewJobInvocationState returns a new job invocation state.
func NewJobInvocationState() *JobInvocationState {
	return &JobInvocationState{
		Output: new(LineWriter),
	}
}

// JobInvocationState is the state object for a job invocation.
type JobInvocationState struct {
	Output *LineWriter
}
