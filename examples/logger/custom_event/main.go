package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/blend/go-sdk/logger"
)

// NewCustomEvent returns a new custom event.
// It is helpful to use a constructor syou can initialize the event meta.
func NewCustomEvent(userID, sessionID, context string) *CustomEvent {
	return &CustomEvent{
		EventMeta: logger.NewEventMeta("custom_event"),
		UserID:    userID,
		SessionID: sessionID,
		Context:   context,
	}
}

// CustomEvent is a custom logger event.
type CustomEvent struct {
	*logger.EventMeta // this embeds a bunch of common fields into our event.

	UserID    string // something domain specific
	SessionID string // something domain specific
	Context   string // something domain specific
}

// WriteText implements logger.TextWritable.
// It is optional, but very much encouraged.
// It takes a formatter and a buffer reference that you push data into.
// This lets the logger re-use buffers.
func (ce CustomEvent) WriteText(tf logger.TextFormatter, wr io.Writer) {
	io.WriteString(wr, ce.UserID)
	io.WriteString(wr, logger.Space)
	io.WriteString(wr, ce.SessionID)
	io.WriteString(wr, logger.Space)
	io.WriteString(wr, ce.Context)
}

// MarshalJSON implements json.Marshaler.
// It is a function that returns just the custom fields on our object as a map,
// to be serialized with the rest of the fields.
func (ce CustomEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(logger.MergeDecomposed(ce.EventMeta.Decompose(), map[string]interface{}{
		"userID":    ce.UserID,
		"sessionID": ce.SessionID,
		"context":   ce.Context,
	}))
}

// NewCustomEventListener returns a type shim for the logger.
func NewCustomEventListener(listener func(context.Context, *CustomEvent)) logger.Listener {
	return func(ctx context.Context, e logger.Event) {
		listener(ctx, e.(*CustomEvent))
	}
}

func main() {
	// make a text logger.
	text := logger.All(logger.OptText())

	// make a json logger
	js := logger.All(logger.OptJSON())

	ctx := context.Background()

	event := NewCustomEvent("bailey", "session0", "Console Demo")

	text.Trigger(ctx, event)
	text.Write(ctx, event)
	js.Trigger(ctx, event)
	js.Write(ctx, event)

	listener := logger.All()
	listener.Listen("custom_event", "demo", NewCustomEventListener(func(_ context.Context, ce *CustomEvent) {
		fmt.Println("listener got event")
	}))
	listener.SyncTrigger(ctx, event)
}
