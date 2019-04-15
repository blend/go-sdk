package logger

import (
	"context"
	"time"
)

// Event is an interface representing methods necessary to trigger listeners.
type Event interface {
	GetFlag() string
	GetTimestamp() time.Time
}

// EventWithContext is an event with a context.
type EventWithContext struct {
	context.Context
	Event
}

// MarshalEvent marshals an object as a logger event.
func MarshalEvent(obj interface{}) (Event, bool) {
	typed, isTyped := obj.(Event)
	return typed, isTyped
}
