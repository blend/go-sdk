package logger

import (
	"context"
	"io"
	"time"
)

// these are compile time assertions
var (
	_ Event = (*MessageEvent)(nil)
)

// NewMessageEvent returns a new message event.
func NewMessageEvent(flag, text string, options ...MessageEventOption) MessageEvent {
	me := MessageEvent{
		Flag: flag,
		Text: text,
	}
	for _, opt := range options {
		opt(&me)
	}
	return me
}

// NewMessageEventListener returns a new message event listener.
func NewMessageEventListener(listener func(context.Context, MessageEvent)) Listener {
	return func(ctx context.Context, e Event) {
		if typed, isTyped := e.(MessageEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// MessageEventOption mutates a message event.
type MessageEventOption func(*MessageEvent)

// OptMessageFlag sets a field on a message event.
func OptMessageFlag(flag string) MessageEventOption {
	return func(me *MessageEvent) { me.Flag = flag }
}

// OptMessageText sets a field on a message event.
func OptMessageText(text string) MessageEventOption {
	return func(me *MessageEvent) { me.Text = text }
}

// OptMessageElapsed sets a field on a message event.
func OptMessageElapsed(elapsed time.Duration) MessageEventOption {
	return func(me *MessageEvent) { me.Elapsed = elapsed }
}

// MessageEvent is a common type of message.
type MessageEvent struct {
	Flag    string
	Text    string
	Elapsed time.Duration
}

// GetFlag implements Event.
func (e MessageEvent) GetFlag() string { return e.Flag }

// WriteText implements TextWritable.
func (e MessageEvent) WriteText(formatter TextFormatter, output io.Writer) {
	io.WriteString(output, e.Text)
	if e.Elapsed > 0 {
		io.WriteString(output, Space)
		io.WriteString(output, "("+e.Elapsed.String()+")")
	}
}

// Decompose implements json.Marshaler.
func (e MessageEvent) Decompose() map[string]interface{} {
	if e.Elapsed > 0 {
		return map[string]interface{}{
			FieldText:    e.Text,
			FieldElapsed: e.Elapsed,
		}
	}
	return map[string]interface{}{
		FieldText: e.Text,
	}
}
