package logger

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/blend/go-sdk/ex"
)

// these are compile time assertions
var (
	_ Event        = (*ErrorEvent)(nil)
	_ TextWritable = (*ErrorEvent)(nil)
	_ JSONWritable = (*ErrorEvent)(nil)
)

// NewErrorEvent returns a new error event.
func NewErrorEvent(flag string, err error, options ...ErrorEventOption) ErrorEvent {
	ee := ErrorEvent{
		Flag: flag,
		Err:  err,
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

// ErrorEventOption is an option for error events.
type ErrorEventOption = func(*ErrorEvent)

// OptErrorEventRequest sets the request on an error event.
func OptErrorEventRequest(req *http.Request) ErrorEventOption {
	return func(ee *ErrorEvent) {
		ee.Request = req
	}
}

// ErrorEvent is an event that wraps an error.
type ErrorEvent struct {
	Flag    string
	Err     error
	Request *http.Request
}

// GetFlag implements Event.
func (ee ErrorEvent) GetFlag() string { return ee.Flag }

// WriteText writes the text version of an error.
func (ee ErrorEvent) WriteText(formatter TextFormatter, output io.Writer) {
	if ee.Err != nil {
		if typed, ok := ee.Err.(*ex.Ex); ok {
			io.WriteString(output, typed.String())
		} else {
			io.WriteString(output, ee.Err.Error())
		}
	}
}

// Decompose implements JSONWritable.
func (ee ErrorEvent) Decompose() map[string]interface{} {
	if ee.Err == nil {
		return nil
	}

	if _, ok := ee.Err.(json.Marshaler); ok {
		return map[string]interface{}{
			"err": ee.Err,
		}
	}
	return map[string]interface{}{
		"err": ee.Err.Error(),
	}
}
