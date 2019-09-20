package logger

import "context"

// EventWithContext is an event with the context it was triggered with.
type EventWithContext struct {
	context.Context
	Event
}
