package logger

import "context"

type skipTriggerKey struct{}

type skipWriteKey struct{}

// WithSkipTrigger sets the context to skip logger listener triggers.
// The event will still be written unless you also use `WithSkipWrite`.
func WithSkipTrigger(ctx context.Context) context.Context {
	return context.WithValue(ctx, skipTriggerKey{}, true)
}

// WithSkipWrite sets the context to skip writing the event to the output stream.
// The event will still trigger listeners unless you also use `WithSkipTrigger`.
func WithSkipWrite(ctx context.Context) context.Context {
	return context.WithValue(ctx, skipWriteKey{}, true)
}

// IsSkipTrigger returns if we should skip triggering logger listeners for a context.
func IsSkipTrigger(ctx context.Context) bool {
	if v := ctx.Value(skipTriggerKey{}); v != nil {
		return true
	}
	return false
}

// IsSkipWrite returns if we should skip writing to the event stream for a context.
func IsSkipWrite(ctx context.Context) bool {
	if v := ctx.Value(skipWriteKey{}); v != nil {
		return true
	}
	return false
}
