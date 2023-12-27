/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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
// A "Shim Writer" is meant to bridge situations where you need to pass
// an io.Writer to a given function, and want that function to write to a logger.
// I.e. you can set `cmd.Stdout = NewShimWriter(log)` to have a
// shell command write to a logger for standard out.
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

// OptShimWriterEventProvider sets the event provider for the shim writer.
func OptShimWriterEventProvider(provider func([]byte) Event) ShimWriterOption {
	return func(sw *ShimWriter) { sw.EventProvider = provider }
}

// OptShimWriterContext sets the context for a given shim writer.
func OptShimWriterContext(ctx context.Context) ShimWriterOption {
	return func(sw *ShimWriter) { sw.Context = ctx }
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
	ctx := sw.Context
	if ctx == nil {
		ctx = context.Background()
	}
	sw.Log.TriggerContext(ctx, sw.EventProvider(contents))
	count = len(contents)
	return
}
