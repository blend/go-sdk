package cron

import (
	"context"

	"github.com/blend/go-sdk/uuid"
)

type jobInvocationKey struct{}

// NewJobInvocationID returns a new pseudo-unique job invocation identifier.
func NewJobInvocationID() string {
	return uuid.V4().String()
}

// WithJobInvocation adds a job invocation to a context as a value.
func WithJobInvocation(ctx context.Context, ji *JobInvocation) context.Context {
	return context.WithValue(ctx, jobInvocationKey{}, ji)
}

// GetJobInvocation returns the job invocation ID from a context.
func GetJobInvocation(ctx context.Context) *JobInvocation {
	if ctx == nil {
		return nil
	}
	if ji, ok := ctx.Value(jobInvocationKey{}).(*JobInvocation); ok {
		return ji
	}
	return nil
}
