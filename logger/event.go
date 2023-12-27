/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package logger

import (
	"context"
	"io"
	"time"
)

// Event is an interface representing methods necessary to trigger listeners.
type Event interface {
	GetFlag() string
}

// TimestampProvider is a type that provides a timestamp.
type TimestampProvider interface {
	GetTimestamp() time.Time
}

// TextWritable is an event that can be written.
type TextWritable interface {
	WriteText(TextFormatter, io.Writer)
}

// JSONWritable is a type that implements a decompose method.
// This is used by the json serializer.
type JSONWritable interface {
	Decompose() map[string]interface{}
}

// GetEventTimestamp returns the timestamp for an event.
// It first checks if the event implements timestamp provider, if it does it returns that value.
// Then it checks if there is a timestamp on the context, if there is one it returns that value.
// Then it checks if there is a triggered timestamp on the context, if there is one it returns that value.
// Then it generates a new timestamp in utc.
func GetEventTimestamp(ctx context.Context, e Event) time.Time {
	if typed, ok := e.(TimestampProvider); ok {
		return typed.GetTimestamp()
	}
	if timestamp := GetTimestamp(ctx); !timestamp.IsZero() {
		return timestamp
	}
	if timestamp := GetTriggerTimestamp(ctx); !timestamp.IsZero() {
		return timestamp
	}
	return time.Now().UTC()
}
