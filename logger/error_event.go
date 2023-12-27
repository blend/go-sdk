/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// these are compile time assertions
var (
	_	Event		= (*ErrorEvent)(nil)
	_	TextWritable	= (*ErrorEvent)(nil)
	_	JSONWritable	= (*ErrorEvent)(nil)
)

// NewErrorEvent returns a new error event.
func NewErrorEvent(flag string, err error, options ...ErrorEventOption) ErrorEvent {
	ee := ErrorEvent{
		Flag:	flag,
		Err:	err,
	}
	for _, opt := range options {
		opt(&ee)
	}
	return ee
}

// NewErrorEventListener returns a new error event listener.
func NewErrorEventListener(listener func(context.Context, ErrorEvent)) Listener {
	return func(ctx context.Context, e Event) {
		if typed, isTyped := e.(ErrorEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// NewScopedErrorEventListener returns a new error event listener that listens
// to specified scopes
func NewScopedErrorEventListener(listener func(context.Context, ErrorEvent), scopes *Scopes) Listener {
	return func(ctx context.Context, e Event) {
		if typed, isTyped := e.(ErrorEvent); isTyped {
			if scopes.IsEnabled(GetPath(ctx)...) {
				listener(ctx, typed)
			}
		}
	}
}

// NewErrorEventFilter returns a new error event filter.
func NewErrorEventFilter(filter func(context.Context, ErrorEvent) (ErrorEvent, bool)) Filter {
	return func(ctx context.Context, e Event) (Event, bool) {
		if typed, isTyped := e.(ErrorEvent); isTyped {
			return filter(ctx, typed)
		}
		return e, false
	}
}

// ErrorEventOption is an option for error events.
type ErrorEventOption = func(*ErrorEvent)

// OptErrorEventState sets the state on an error event.
func OptErrorEventState(state interface{}) ErrorEventOption {
	return func(ee *ErrorEvent) {
		ee.State = state
	}
}

// OptErrorEventRestricted sets the restricted field on an error event.
func OptErrorEventRestricted(restricted bool) ErrorEventOption {
	return func(ee *ErrorEvent) {
		ee.Restricted = restricted
	}
}

// ErrorEvent is an event that wraps an error.
type ErrorEvent struct {
	Flag		string
	Err		error
	State		interface{}
	Restricted	bool
}

// GetFlag implements Event.
func (ee ErrorEvent) GetFlag() string	{ return ee.Flag }

// WriteText writes the text version of an error.
func (ee ErrorEvent) WriteText(formatter TextFormatter, output io.Writer) {
	if ee.Err != nil {
		fmt.Fprintf(output, "%+v", ee.Err)
	}
}

// Decompose implements JSONWritable.
func (ee ErrorEvent) Decompose() map[string]interface{} {
	if ee.Err == nil {
		return nil
	}

	if _, ok := ee.Err.(json.Marshaler); ok {
		return map[string]interface{}{
			"err":			ee.Err,
			FieldRestricted:	ee.Restricted,
		}
	}
	return map[string]interface{}{
		"err":			ee.Err.Error(),
		FieldRestricted:	ee.Restricted,
	}
}
