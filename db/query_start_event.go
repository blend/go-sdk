/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package db

import (
	"context"
	"fmt"
	"io"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stringutil"
)

// Logger flags
const (
	QueryStartFlag = "db.query.start"
)

// these are compile time assertions
var (
	_	logger.Event		= (*QueryEvent)(nil)
	_	logger.TextWritable	= (*QueryEvent)(nil)
	_	logger.JSONWritable	= (*QueryEvent)(nil)
)

// these are compile time assertions
var (
	_	logger.Event		= (*QueryEvent)(nil)
	_	logger.TextWritable	= (*QueryEvent)(nil)
	_	logger.JSONWritable	= (*QueryEvent)(nil)
)

// NewQueryStartEvent creates a new query start event.
func NewQueryStartEvent(body string, options ...QueryStartEventOption) QueryStartEvent {
	qse := QueryStartEvent{
		Body: body,
	}
	for _, opt := range options {
		opt(&qse)
	}
	return qse
}

// NewQueryStartEventListener returns a new listener for query events.
func NewQueryStartEventListener(listener func(context.Context, QueryStartEvent)) logger.Listener {
	return func(ctx context.Context, e logger.Event) {
		if typed, isTyped := e.(QueryStartEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// NewQueryStartEventFilter returns a new query event filter.
func NewQueryStartEventFilter(filter func(context.Context, QueryStartEvent) (QueryStartEvent, bool)) logger.Filter {
	return func(ctx context.Context, e logger.Event) (logger.Event, bool) {
		if typed, isTyped := e.(QueryStartEvent); isTyped {
			return filter(ctx, typed)
		}
		return e, false
	}
}

// QueryStartEventOption mutates a query start event.
type QueryStartEventOption func(*QueryStartEvent)

// OptQueryStartEventBody sets a field on the query event.
func OptQueryStartEventBody(value string) QueryStartEventOption {
	return func(e *QueryStartEvent) { e.Body = value }
}

// OptQueryStartEventDatabase sets a field on the query event.
func OptQueryStartEventDatabase(value string) QueryStartEventOption {
	return func(e *QueryStartEvent) { e.Database = value }
}

// OptQueryStartEventEngine sets a field on the query event.
func OptQueryStartEventEngine(value string) QueryStartEventOption {
	return func(e *QueryStartEvent) { e.Engine = value }
}

// OptQueryStartEventUsername sets a field on the query event.
func OptQueryStartEventUsername(value string) QueryStartEventOption {
	return func(e *QueryStartEvent) { e.Username = value }
}

// OptQueryStartEventLabel sets a field on the query event.
func OptQueryStartEventLabel(label string) QueryStartEventOption {
	return func(e *QueryStartEvent) { e.Label = label }
}

// QueryStartEvent represents the start of a database query.
type QueryStartEvent struct {
	Database	string
	Engine		string
	Username	string
	Label		string
	Body		string
}

// GetFlag implements Event.
func (e QueryStartEvent) GetFlag() string	{ return QueryStartFlag }

// WriteText writes the event text to the output.
func (e QueryStartEvent) WriteText(tf logger.TextFormatter, wr io.Writer) {
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
}

// Decompose implements JSONWritable.
func (e QueryStartEvent) Decompose() map[string]interface{} {
	return map[string]interface{}{
		"engine":	e.Engine,
		"database":	e.Database,
		"username":	e.Username,
		"label":	e.Label,
		"body":		e.Body,
	}
}
