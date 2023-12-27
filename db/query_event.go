/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

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
	_	logger.Event		= (*QueryEvent)(nil)
	_	logger.TextWritable	= (*QueryEvent)(nil)
	_	logger.JSONWritable	= (*QueryEvent)(nil)
)

// NewQueryEvent creates a new query event.
func NewQueryEvent(body string, elapsed time.Duration, options ...QueryEventOption) QueryEvent {
	qe := QueryEvent{
		Body:		body,
		Elapsed:	elapsed,
	}
	for _, opt := range options {
		opt(&qe)
	}
	return qe
}

// NewQueryEventListener returns a new listener for query events.
func NewQueryEventListener(listener func(context.Context, QueryEvent)) logger.Listener {
	return func(ctx context.Context, e logger.Event) {
		if typed, isTyped := e.(QueryEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// NewQueryEventFilter returns a new query event filter.
func NewQueryEventFilter(filter func(context.Context, QueryEvent) (QueryEvent, bool)) logger.Filter {
	return func(ctx context.Context, e logger.Event) (logger.Event, bool) {
		if typed, isTyped := e.(QueryEvent); isTyped {
			return filter(ctx, typed)
		}
		return e, false
	}
}

// QueryEventOption mutates a query event.
type QueryEventOption func(*QueryEvent)

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

// OptQueryEventLabel sets a field on the query event.
func OptQueryEventLabel(label string) QueryEventOption {
	return func(e *QueryEvent) { e.Label = label }
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
	Database	string
	Engine		string
	Username	string
	Label		string
	Body		string
	Elapsed		time.Duration
	Err		error
}

// GetFlag implements Event.
func (e QueryEvent) GetFlag() string	{ return QueryFlag }

// WriteText writes the event text to the output.
func (e QueryEvent) WriteText(tf logger.TextFormatter, wr io.Writer) {
	fmt.Fprint(wr, "[")
	if len(e.Engine) > 0 {
		fmt.Fprint(wr, tf.Colorize(e.Engine, ansi.ColorLightWhite))
		fmt.Fprint(wr, logger.Space)
	}
	if len(e.Username) > 0 {
		fmt.Fprint(wr, tf.Colorize(e.Username, ansi.ColorLightWhite))
		fmt.Fprint(wr, "@")
	}
	fmt.Fprint(wr, tf.Colorize(e.Database, ansi.ColorLightWhite))
	fmt.Fprint(wr, "]")

	if len(e.Label) > 0 {
		fmt.Fprint(wr, logger.Space)
		fmt.Fprintf(wr, "[%s]", tf.Colorize(e.Label, ansi.ColorLightWhite))
	}

	if len(e.Body) > 0 {
		fmt.Fprint(wr, logger.Space)
		fmt.Fprint(wr, stringutil.CompressSpace(e.Body))
	}

	fmt.Fprint(wr, logger.Space)
	fmt.Fprint(wr, e.Elapsed.String())

	if e.Err != nil {
		fmt.Fprint(wr, logger.Space)
		fmt.Fprint(wr, tf.Colorize("failed", ansi.ColorRed))
	}
}

// Decompose implements JSONWritable.
func (e QueryEvent) Decompose() map[string]interface{} {
	return map[string]interface{}{
		"engine":	e.Engine,
		"database":	e.Database,
		"username":	e.Username,
		"label":	e.Label,
		"body":		e.Body,
		"err":		e.Err,
		"elapsed":	timeutil.Milliseconds(e.Elapsed),
	}
}
