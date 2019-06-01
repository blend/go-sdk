package logger

import (
	"context"
	"encoding/json"
	"io"
	"time"
)

// these are compile time assertions
var (
	_ Event = (*MessageEvent)(nil)
)

// NewMessageEvent returns a new message event.
func NewMessageEvent(flag, message string, options ...MessageEventOption) *MessageEvent {
	me := MessageEvent{
		EventMeta: NewEventMeta(flag),
		Message:   message,
	}
	for _, opt := range options {
		opt(&me)
	}
	return &me
}

// NewMessageEventListener returns a new message event listener.
func NewMessageEventListener(listener func(context.Context, *MessageEvent)) Listener {
	return func(ctx context.Context, e Event) {
		if typed, isTyped := e.(*MessageEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// MessageEventOption mutates a message event.
type MessageEventOption func(*MessageEvent)

// OptMessageMeta sets meta options.
func OptMessageMeta(options ...EventMetaOption) MessageEventOption {
	return func(me *MessageEvent) {
		for _, opt := range options {
			opt(me.EventMeta)
		}
	}
}

// OptMessage sets a field on a message event.
// Code style note; `OptMessageMessage` stutters, so it's been shortened.
func OptMessage(message string) MessageEventOption {
	return func(me *MessageEvent) { me.Message = message }
}

// OptMessageElapsed sets a field on a message event.
func OptMessageElapsed(elapsed time.Duration) MessageEventOption {
	return func(me *MessageEvent) { me.Elapsed = elapsed }
}

// MessageEvent is a common type of message.
type MessageEvent struct {
	*EventMeta `json:",inline"`
	Message    string        `json:"message"`
	Elapsed    time.Duration `json:"elapsed"`
}

// WriteText implements TextWritable.
func (e *MessageEvent) WriteText(formatter TextFormatter, output io.Writer) {
	io.WriteString(output, e.Message)
	if e.Elapsed > 0 {
		io.WriteString(output, Space)
		io.WriteString(output, "("+e.Elapsed.String()+")")
	}
}

// MarshalJSON implements json.Marshaler.
func (e MessageEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(MergeDecomposed(e.EventMeta.Decompose(), map[string]interface{}{
		"message": e.Message,
	}))
}
