package main

import (
	"context"
	"fmt"
	"io"

	"github.com/blend/go-sdk/logger"
)

// NewCustomEvent returns a new custom event.
// It is helpful to use a constructor syou can initialize the event meta.
func NewCustomEvent(userID, sessionID, context string) CustomEvent {
	return CustomEvent{
		UserID:    userID,
		SessionID: sessionID,
		Context:   context,
	}
}

// CustomEvent is a custom logger event.
type CustomEvent struct {
	UserID    string // something domain specific
	SessionID string // something domain specific
	Context   string // something domain specific
}

// GetFlag implements logger.Event.
func (ce CustomEvent) GetFlag() string { return "custom_event" }

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

// Decompose implements logger.JSONWritable.
// It is a function that returns just the custom fields on our object as a map,
// to be serialized with the rest of the fields.
func (ce CustomEvent) Decompose() map[string]interface{} {
	return map[string]interface{}{
		"userID":    ce.UserID,
		"sessionID": ce.SessionID,
		"context":   ce.Context,
	}
}

// NewCustomEventListener returns a type shim for the logger.
func NewCustomEventListener(listener func(context.Context, CustomEvent)) logger.Listener {
	return func(ctx context.Context, e logger.Event) {
		listener(ctx, e.(CustomEvent))
	}
}

func main() {
	// make a text logger.
	text := logger.All(logger.OptText())

	// make a json logger
	js := logger.All(logger.OptJSON())

	ctx := context.Background()

	event := NewCustomEvent("bailey", "session0", "Console Demo")

	text.TriggerContext(ctx, event)
	text.Write(ctx, event)
	js.TriggerContext(ctx, event)
	js.Write(ctx, event)

	done := make(chan struct{})
	listener := logger.All()
	listener.Listen("custom_event", "demo", NewCustomEventListener(func(_ context.Context, ce CustomEvent) {
		fmt.Println("listener got event")
		close(done)
	}))
	listener.TriggerContext(ctx, event)
	<-done
}
