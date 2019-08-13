package logger

import (
	"context"
	"encoding/json"
	"io"

	"github.com/blend/go-sdk/ex"
)

// these are compile time assertions
var (
	_ Event          = (*ErrorEvent)(nil)
	_ TextWritable   = (*ErrorEvent)(nil)
	_ json.Marshaler = (*ErrorEvent)(nil)
)

// NewErrorEvent returns a new error event.
func NewErrorEvent(flag string, err error, options ...ErrorEventOption) *ErrorEvent {
	ee := &ErrorEvent{
		EventMeta: NewEventMeta(
			flag,
		),
		Err: err,
	}
	for _, option := range options {
		option(ee)
	}
	return ee
}

// NewErrorEventListener returns a new error event listener.
func NewErrorEventListener(listener func(context.Context, *ErrorEvent)) Listener {
	return func(ctx context.Context, e Event) {
		if typed, isTyped := e.(*ErrorEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// ErrorEventOption is an option for ErrorEvents.
type ErrorEventOption func(*ErrorEvent)

// OptErrorEventMetaOptions sets the event meta options.
func OptErrorEventMetaOptions(options ...EventMetaOption) ErrorEventOption {
	return func(e *ErrorEvent) {
		for _, option := range options {
			option(e.EventMeta)
		}
	}
}

// OptErrorEventErr sets the error on the error event.
func OptErrorEventErr(err error) ErrorEventOption {
	return func(e *ErrorEvent) { e.Err = err }
}

// OptErrorEventState sets the state on the error event.
func OptErrorEventState(state interface{}) ErrorEventOption {
	return func(e *ErrorEvent) { e.State = state }
}

// ErrorEvent is an event that wraps an error.
type ErrorEvent struct {
	*EventMeta
	Err   error
	State interface{}
}

// WriteText writes the text version of an error.
func (e ErrorEvent) WriteText(formatter TextFormatter, output io.Writer) {
	if e.Err != nil {
		if typed, ok := e.Err.(*ex.Ex); ok {
			io.WriteString(output, typed.String())
		} else {
			io.WriteString(output, e.Err.Error())
		}
	}
}

// MarshalJSON implements json.Marshaler.
func (e ErrorEvent) MarshalJSON() ([]byte, error) {
	if _, ok := e.Err.(json.Marshaler); ok {
		return json.Marshal(MergeDecomposed(e.EventMeta.Decompose(), map[string]interface{}{
			"err":   e.Err,
			"state": e.State,
		}))
	}
	return json.Marshal(MergeDecomposed(e.EventMeta.Decompose(), map[string]interface{}{
		"err":   e.Err.Error(),
		"state": e.State,
	}))
}
