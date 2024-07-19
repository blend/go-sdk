/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"context"
	"time"
)

type loggerKey struct{}

// WithLogger adds the logger to a context.
func WithLogger(ctx context.Context, log Log) context.Context {
	return context.WithValue(ctx, loggerKey{}, log)
}

// GetLogger gets a logger off a context.
func GetLogger(ctx context.Context) Log {
	if value := ctx.Value(loggerKey{}); value != nil {
		if typed, ok := value.(Log); ok {
			return typed
		}
	}
	return nil
}

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

// WithPath returns a new context with a given path segment(s).
//
// NOTE: This overwrites any _existing_ path context values.
//
// If you want to _append_ path segments, use `logger.WithPathAppend(...)`.
func WithPath(ctx context.Context, path ...string) context.Context {
	return context.WithValue(ctx, pathKey{}, path)
}

// WithPathAppend appends a given path segment to a context.
func WithPathAppend(ctx context.Context, path ...string) context.Context {
	return context.WithValue(ctx, pathKey{}, append(GetPath(ctx), path...))
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

// WithSetLabels sets the labels on a context, overwriting any existing labels.
func WithSetLabels(ctx context.Context, labels Labels) context.Context {
	return context.WithValue(ctx, labelsKey{}, labels)
}

// WithLabels returns a new context with a given additional labels.
//
// Any labels already on the context will be added to the new set
// with the provided set being layered on top of those.
func WithLabels(ctx context.Context, labels Labels) context.Context {
	// it is critical here that we copy the labels
	// and not mutate the existing map
	combinedLabels := GetLabels(ctx)
	for key, value := range labels {
		combinedLabels[key] = value
	}
	return context.WithValue(ctx, labelsKey{}, combinedLabels)
}

// WithLabel returns a new context with a given additional label.
//
// Any labels already on the context will be added to the new set
// with the provided set being layered on top of those.
func WithLabel(ctx context.Context, key, value string) context.Context {
	newLabels := make(Labels)

	// it is critical here that we copy the labels
	// and not mutate the existing map.
	existing := GetLabels(ctx)
	for key, value := range existing {
		newLabels[key] = value
	}

	// we assign after as the new value should overwrite
	newLabels[key] = value
	return context.WithValue(ctx, labelsKey{}, newLabels)
}

// GetLabels gets labels off a context.
//
// It will return a copy of the labels, preventing map races.
func GetLabels(ctx context.Context) Labels {
	if raw := ctx.Value(labelsKey{}); raw != nil {
		if typed, ok := raw.(Labels); ok {
			// create a copy
			output := make(Labels)
			for key, value := range typed {
				output[key] = value
			}
			return output
		}
	}
	return make(Labels)
}

// GetLabel gets a label off a context by key.
func GetLabel(ctx context.Context, key string) (string, bool) {
	if raw := ctx.Value(labelsKey{}); raw != nil {
		if typed, ok := raw.(Labels); ok {
			if value, found := typed[key]; found {
				return value, true
			}
		}
	}
	return "", false
}

type annotationsKey struct{}

// WithAnnotations returns a new context with a given additional annotations.
func WithAnnotations(ctx context.Context, annotations Annotations) context.Context {
	return context.WithValue(ctx, annotationsKey{}, annotations)
}

// WithAnnotation returns a new context with a given additional annotation.
func WithAnnotation(ctx context.Context, key string, value interface{}) context.Context {
	existing := GetAnnotations(ctx)
	existing[key] = value
	return context.WithValue(ctx, annotationsKey{}, existing)
}

// GetAnnotations gets annotations off a context.
func GetAnnotations(ctx context.Context) Annotations {
	if raw := ctx.Value(annotationsKey{}); raw != nil {
		if typed, ok := raw.(Annotations); ok {
			// create a copy
			output := make(Annotations)
			for key, value := range typed {
				output[key] = value
			}
			return output
		}
	}
	return make(Annotations)
}

type restrictedKey struct{}

// WithRestricted returns a new context with the restricted value.
func WithRestricted(ctx context.Context, restricted bool) context.Context {
	return context.WithValue(ctx, restrictedKey{}, restricted)
}

// GetRestricted gets restricted value from a context.
func GetRestricted(ctx context.Context) bool {
	if raw := ctx.Value(restrictedKey{}); raw != nil {
		if typed, ok := raw.(bool); ok {
			return typed
		}
	}
	return false
}

type skipTriggerKey struct{}

type skipWriteKey struct{}

// WithSkipTrigger sets the context to skip logger listener triggers.
// The event will still be written unless you also use `WithSkipWrite`.
func WithSkipTrigger(ctx context.Context, skipTrigger bool) context.Context {
	return context.WithValue(ctx, skipTriggerKey{}, skipTrigger)
}

// WithSkipWrite sets the context to skip writing the event to the output stream.
// The event will still trigger listeners unless you also use `WithSkipTrigger`.
func WithSkipWrite(ctx context.Context, skipWrite bool) context.Context {
	return context.WithValue(ctx, skipWriteKey{}, skipWrite)
}

// IsSkipTrigger returns if we should skip triggering logger listeners for a context.
func IsSkipTrigger(ctx context.Context) bool {
	if v := ctx.Value(skipTriggerKey{}); v != nil {
		if typed, ok := v.(bool); ok {
			return typed
		}
	}
	return false
}

// IsSkipWrite returns if we should skip writing to the event stream for a context.
func IsSkipWrite(ctx context.Context) bool {
	if v := ctx.Value(skipWriteKey{}); v != nil {
		if typed, ok := v.(bool); ok {
			return typed
		}
	}
	return false
}
