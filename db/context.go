package db

import (
	"context"
)

type connectionKey struct{}

// WithConnection adds a given connection to the context.
func WithConnection(ctx context.Context, conn *Connection) context.Context {
	return context.WithValue(ctx, connectionKey{}, conn)
}

// GetConnection adds a given connection to the context.
func GetConnection(ctx context.Context) *Connection {
	if value := ctx.Value(connectionKey{}); value != nil {
		if typed, ok := value.(*Connection); ok {
			return typed
		}
	}
	return nil
}

type skipQueryLogging struct{}

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
