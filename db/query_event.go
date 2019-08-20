package db

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stringutil"
	"github.com/blend/go-sdk/timeutil"
)

// Logger flags
const (
	QueryFlag = "db.query"
)

// these are compile time assertions
var (
	_ logger.Event        = (*QueryEvent)(nil)
	_ logger.TextWritable = (*QueryEvent)(nil)
	_ logger.JSONWritable = (*QueryEvent)(nil)
)

// NewQueryEvent creates a new query event.
func NewQueryEvent(body string, elapsed time.Duration, options ...QueryEventOption) QueryEvent {
	qe := QueryEvent{
		Body:    body,
		Elapsed: elapsed,
	}
	for _, opt := range options {
		opt(&qe)
	}
	return qe
}

// NewQueryEventListener returns a new listener for spiffy events.
func NewQueryEventListener(listener func(context.Context, QueryEvent)) logger.Listener {
	return func(ctx context.Context, e logger.Event) {
		if typed, isTyped := e.(QueryEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// QueryEventOption mutates a query event.
type QueryEventOption func(*QueryEvent)

// OptQueryBody sets a field on the query event.
func OptQueryBody(value string) QueryEventOption {
	return func(e *QueryEvent) { e.Body = value }
}

// OptQueryDatabase sets a field on the query event.
func OptQueryDatabase(value string) QueryEventOption {
	return func(e *QueryEvent) { e.Database = value }
}

// OptQueryEngine sets a field on the query event.
func OptQueryEngine(value string) QueryEventOption {
	return func(e *QueryEvent) { e.Engine = value }
}

// OptQueryUsername sets a field on the query event.
func OptQueryUsername(value string) QueryEventOption {
	return func(e *QueryEvent) { e.Username = value }
}

// OptQueryLabel sets a field on the query event.
func OptQueryLabel(value string) QueryEventOption {
	return func(e *QueryEvent) { e.QueryLabel = value }
}

// OptQueryElapsed sets a field on the query event.
func OptQueryElapsed(value time.Duration) QueryEventOption {
	return func(e *QueryEvent) { e.Elapsed = value }
}

// OptQueryErr sets a field on the query event.
func OptQueryErr(value error) QueryEventOption {
	return func(e *QueryEvent) { e.Err = value }
}

// QueryEvent represents a database query.
type QueryEvent struct {
	Database   string
	Engine     string
	Username   string
	QueryLabel string
	Body       string
	Elapsed    time.Duration
	Err        error
}

// GetFlag implements Event.
func (e QueryEvent) GetFlag() string { return QueryFlag }

// WriteText writes the event text to the output.
func (e QueryEvent) WriteText(tf logger.TextFormatter, wr io.Writer) {
	io.WriteString(wr, "[")
	if len(e.Engine) > 0 {
		io.WriteString(wr, tf.Colorize(e.Engine, ansi.ColorLightWhite))
		io.WriteString(wr, logger.Space)
	}
	if len(e.Username) > 0 {
		io.WriteString(wr, tf.Colorize(e.Username, ansi.ColorLightWhite))
		io.WriteString(wr, "@")
	}
	io.WriteString(wr, tf.Colorize(e.Database, ansi.ColorLightWhite))
	io.WriteString(wr, "]")

	if len(e.QueryLabel) > 0 {
		io.WriteString(wr, logger.Space)
		io.WriteString(wr, fmt.Sprintf("[%s]", tf.Colorize(e.QueryLabel, ansi.ColorLightWhite)))
	}

	io.WriteString(wr, logger.Space)
	io.WriteString(wr, e.Elapsed.String())

	if e.Err != nil {
		io.WriteString(wr, logger.Space)
		io.WriteString(wr, tf.Colorize("failed", ansi.ColorRed))
	}

	if len(e.Body) > 0 {
		io.WriteString(wr, logger.Space)
		io.WriteString(wr, stringutil.CompressSpace(e.Body))
	}
}

// Decompose implements JSONWritable.
func (e QueryEvent) Decompose() map[string]interface{} {
	return map[string]interface{}{
		"engine":     e.Engine,
		"database":   e.Database,
		"username":   e.Username,
		"queryLabel": e.QueryLabel,
		"body":       e.Body,
		"err":        e.Err,
		"elapsed":    timeutil.Milliseconds(e.Elapsed),
	}
}
