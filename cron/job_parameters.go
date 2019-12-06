package cron

import (
	"context"
)

// JobParameters is a loose association to map[string]string.
type JobParameters = map[string]string

type contextKeyJobParameters struct{}

// WithJobParameters adds job invocation parameter values to a context.
func WithJobParameters(ctx context.Context, values JobParameters) context.Context {
	return context.WithValue(ctx, contextKeyJobParameters{}, values)
}

// GetJobParameters gets parameter values from a given context.
func GetJobParameters(ctx context.Context) JobParameters {
	if value := ctx.Value(contextKeyJobParameters{}); value != nil {
		if typed, ok := value.(JobParameters); ok {
			return typed
		}
	}
	return nil
}
