package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/blend/go-sdk/stringutil"
	"github.com/blend/go-sdk/timeutil"

	"github.com/blend/go-sdk/ansi"
)

// these are compile time assertions
var (
	_ Event          = (*QueryEvent)(nil)
	_ TextWritable   = (*QueryEvent)(nil)
	_ json.Marshaler = (*QueryEvent)(nil)
)

// NewQueryEvent creates a new query event.
func NewQueryEvent(body string, elapsed time.Duration, options ...QueryEventOption) *QueryEvent {
	qe := QueryEvent{
		EventMeta: NewEventMeta(Query),
		Body:      body,
		Elapsed:   elapsed,
	}
	for _, opt := range options {
		opt(&qe)
	}
	return &qe
}

// NewQueryEventListener returns a new listener for spiffy events.
func NewQueryEventListener(listener func(context.Context, *QueryEvent)) Listener {
	return func(ctx context.Context, e Event) {
		if typed, isTyped := e.(*QueryEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// QueryEventOption mutates a query event.
type QueryEventOption func(*QueryEvent)

// OptQueryEventMetaOptions sets options on the event metadata.
func OptQueryEventMetaOptions(options ...EventMetaOption) QueryEventOption {
	return func(ae *QueryEvent) {
		for _, option := range options {
			option(ae.EventMeta)
		}
	}
}

// OptQueryEventBody sets a field on the query event.
func OptQueryEventBody(value string) QueryEventOption {
	return func(e *QueryEvent) { e.Body = value }
}

// OptQueryEventDatabase sets a field on the query event.
func OptQueryEventDatabase(value string) QueryEventOption {
	return func(e *QueryEvent) { e.Database = value }
}

// OptQueryEventEngine sets a field on the query event.
func OptQueryEventEngine(value string) QueryEventOption {
	return func(e *QueryEvent) { e.Engine = value }
}

// OptQueryEventUsername sets a field on the query event.
func OptQueryEventUsername(value string) QueryEventOption {
	return func(e *QueryEvent) { e.Username = value }
}

// OptQueryEventQueryLabel sets a field on the query event.
func OptQueryEventQueryLabel(value string) QueryEventOption {
	return func(e *QueryEvent) { e.QueryLabel = value }
}

// OptQueryEventElapsed sets a field on the query event.
func OptQueryEventElapsed(value time.Duration) QueryEventOption {
	return func(e *QueryEvent) { e.Elapsed = value }
}

// OptQueryEventErr sets a field on the query event.
func OptQueryEventErr(value error) QueryEventOption {
	return func(e *QueryEvent) { e.Err = value }
}

// QueryEvent represents a database query.
type QueryEvent struct {
	*EventMeta

	Database   string
	Engine     string
	Username   string
	QueryLabel string
	Body       string
	Elapsed    time.Duration
	Err        error
}

// WriteText writes the event text to the output.
func (e QueryEvent) WriteText(tf TextFormatter, wr io.Writer) {
	io.WriteString(wr, "[")
	if len(e.Engine) > 0 {
		io.WriteString(wr, tf.Colorize(e.Engine, ansi.ColorLightWhite))
		io.WriteString(wr, Space)
	}
	if len(e.Username) > 0 {
		io.WriteString(wr, tf.Colorize(e.Username, ansi.ColorLightWhite))
		io.WriteString(wr, "@")
	}
	io.WriteString(wr, tf.Colorize(e.Database, ansi.ColorLightWhite))
	io.WriteString(wr, "]")

	if len(e.QueryLabel) > 0 {
		io.WriteString(wr, Space)
		io.WriteString(wr, fmt.Sprintf("[%s]", tf.Colorize(e.QueryLabel, ansi.ColorLightWhite)))
	}

	io.WriteString(wr, Space)
	io.WriteString(wr, e.Elapsed.String())

	if e.Err != nil {
		io.WriteString(wr, Space)
		io.WriteString(wr, tf.Colorize("failed", ansi.ColorRed))
	}

	if len(e.Body) > 0 {
		io.WriteString(wr, Space)
		io.WriteString(wr, stringutil.CompressSpace(e.Body))
	}
}

// MarshalJSON implements json.Marshaler.
func (e QueryEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(MergeDecomposed(e.EventMeta.Decompose(), map[string]interface{}{
		"engine":     e.Engine,
		"database":   e.Database,
		"username":   e.Username,
		"queryLabel": e.QueryLabel,
		"body":       e.Body,
		"err":        e.Err,
		"elapsed":    timeutil.Milliseconds(e.Elapsed),
	}))
}
