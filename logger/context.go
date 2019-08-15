package logger

import (
	"context"
	"time"
)

type timestampKey struct{}

// WithTimestamp returns a new context with a given timestamp value.
func WithTimestamp(ctx context.Context, ts time.Time) context.Context {
	return context.WithValue(ctx, timestampKey{}, ts)
}

// GetTimestamp gets a timestampoff a context.
func GetTimestamp(ctx context.Context) time.Time {
	if raw := ctx.Value(timestampKey{}); raw != nil {
		if typed, ok := raw.(time.Time); ok {
			return typed
		}
	}
	return time.Now().UTC()
}

type scopePathKey struct{}

// WithScopePath returns a new context with a given additional path segment(s).
func WithScopePath(ctx context.Context, path ...string) context.Context {
	return context.WithValue(ctx, scopePathKey{}, append(GetScopePath(ctx), path...))
}

// GetScopePath gets a scope path off a context.
func GetScopePath(ctx context.Context) []string {
	if raw := ctx.Value(scopePathKey{}); raw != nil {
		if typed, ok := raw.([]string); ok {
			return typed
		}
	}
	return nil
}

type fieldsKey struct{}

// WithFields returns a new context with a given additional path segments.
func WithFields(ctx context.Context, fields Fields) context.Context {
	return context.WithValue(ctx, fieldsKey{}, CombineFields(GetFields(ctx), fields))
}

// GetFields gets fields off a context.
func GetFields(ctx context.Context) Fields {
	if raw := ctx.Value(fieldsKey{}); raw != nil {
		if typed, ok := raw.(Fields); ok {
			return typed
		}
	}
	return nil
}
