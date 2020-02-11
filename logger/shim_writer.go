package logger

import (
	"context"
	"fmt"
	"io"
	"strings"
)

var (
	_ io.Writer = (*ShimWriter)(nil)
)

// Constants
const (
	DefaultShimWriterMessageFlag = "shim"
)

// NewShimWriter returns a new shim writer.
func NewShimWriter(log Triggerable, opts ...ShimWriterOption) ShimWriter {
	shim := ShimWriter{
		Context:       context.Background(),
		Log:           log,
		EventProvider: ShimWriterMessageEventProvider(DefaultShimWriterMessageFlag),
	}
	for _, opt := range opts {
		opt(&shim)
	}
	return shim
}

// OptShimWriterContext sets the base context for the shim writer.
func OptShimWriterContext(ctx context.Context) ShimWriterOption {
	return func(sw *ShimWriter) { sw.Context = ctx }
}

// OptShimWriterEventProvider sets the event provider for the shim writer.
func OptShimWriterEventProvider(provider func([]byte) Event) ShimWriterOption {
	return func(sw *ShimWriter) { sw.EventProvider = provider }
}

// ShimWriterMessageEventProvider returns a message event with a given flag
// for a given contents.
func ShimWriterMessageEventProvider(flag string, opts ...MessageEventOption) func([]byte) Event {
	return func(contents []byte) Event {
		return NewMessageEvent(flag, strings.TrimSpace(string(contents)), opts...)
	}
}

// ShimWriterErrorEventProvider returns an error event with a given flag
// for a given contents.
func ShimWriterErrorEventProvider(flag string, opts ...ErrorEventOption) func([]byte) Event {
	return func(contents []byte) Event {
		return NewErrorEvent(flag, fmt.Errorf(strings.TrimSpace(string(contents))), opts...)
	}
}

// ShimWriterOption is a mutator for a shim writer.
type ShimWriterOption func(*ShimWriter)

// ShimWriter is a type that implements io.Writer with
// a logger backend.
type ShimWriter struct {
	Context       context.Context
	Log           Triggerable
	EventProvider func([]byte) Event
}

// Write implements io.Writer.
func (sw ShimWriter) Write(contents []byte) (count int, err error) {
	sw.Log.Trigger(sw.Context, sw.EventProvider(contents))
	count = len(contents)
	return
}
