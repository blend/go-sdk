package db

import "context"

type skipQueryLogging struct {}

// WithSkipQueryLogging sets the context to skip logger listener triggers.
func WithSkipQueryLogging(ctx context.Context) context.Context {
	return context.WithValue(ctx, skipQueryLogging{}, true)
}

// IsSkipQueryLogging returns if we should skip triggering logger listeners for a context.
func IsSkipQueryLogging(ctx context.Context) bool {
	if v := ctx.Value(skipQueryLogging{}); v != nil {
		return true
	}
	return false
}
