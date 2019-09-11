package logger

import (
	"context"
	"time"
)

type triggerTimestampKey struct{}

// WithTriggerTimestamp returns a new context with a given timestamp value.
// It is used by the scope to connote when an event was triggered.
func WithTriggerTimestamp(ctx context.Context, ts time.Time) context.Context {
	return context.WithValue(ctx, triggerTimestampKey{}, ts)
}

// GetTriggerTimestamp gets when an event was triggered off a context.
func GetTriggerTimestamp(ctx context.Context) time.Time {
	if raw := ctx.Value(triggerTimestampKey{}); raw != nil {
		if typed, ok := raw.(time.Time); ok {
			return typed
		}
	}
	return time.Time{}
}

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
	return time.Time{}
}

type pathKey struct{}

// WithPath returns a new context with a given additional path segment(s).
func WithPath(ctx context.Context, path ...string) context.Context {
	return context.WithValue(ctx, pathKey{}, path)
}

// GetPath gets a path off a context.
func GetPath(ctx context.Context) []string {
	if raw := ctx.Value(pathKey{}); raw != nil {
		if typed, ok := raw.([]string); ok {
			return typed
		}
	}
	return nil
}

type labelsKey struct{}

// WithLabels returns a new context with a given additional labels.
func WithLabels(ctx context.Context, labels Labels) context.Context {
	return context.WithValue(ctx, labelsKey{}, labels)
}

// GetLabels gets labels off a context.
func GetLabels(ctx context.Context) Labels {
	if raw := ctx.Value(labelsKey{}); raw != nil {
		if typed, ok := raw.(Labels); ok {
			return typed
		}
	}
	return nil
}

type annotationsKey struct{}

// WithAnnotations returns a new context with a given additional annotations.
func WithAnnotations(ctx context.Context, annotations Annotations) context.Context {
	return context.WithValue(ctx, annotationsKey{}, annotations)
}

// GetAnnotations gets annotations off a context.
func GetAnnotations(ctx context.Context) Annotations {
	if raw := ctx.Value(annotationsKey{}); raw != nil {
		if typed, ok := raw.(Annotations); ok {
			return typed
		}
	}
	return nil
}
