/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package redis

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/timeutil"
)

// Logger flags
const (
	Flag = "redis"
)

// these are compile time assertions
var (
	_	logger.Event		= (*Event)(nil)
	_	logger.TextWritable	= (*Event)(nil)
	_	logger.JSONWritable	= (*Event)(nil)
)

// NewEvent creates a new query event.
func NewEvent(op string, args []string, elapsed time.Duration, options ...EventOption) Event {
	qe := Event{
		Op:		op,
		Args:		args,
		Elapsed:	elapsed,
	}
	for _, opt := range options {
		opt(&qe)
	}
	return qe
}

// NewEventListener returns a new listener for events.
func NewEventListener(listener func(context.Context, Event)) logger.Listener {
	return func(ctx context.Context, e logger.Event) {
		if typed, isTyped := e.(Event); isTyped {
			listener(ctx, typed)
		}
	}
}

// NewEventFilter returns a new event filter.
func NewEventFilter(filter func(context.Context, Event) (Event, bool)) logger.Filter {
	return func(ctx context.Context, e logger.Event) (logger.Event, bool) {
		if typed, isTyped := e.(Event); isTyped {
			return filter(ctx, typed)
		}
		return e, false
	}
}

// EventOption mutates an event.
type EventOption func(*Event)

// OptEventNetwork sets a field on the event.
func OptEventNetwork(value string) EventOption {
	return func(e *Event) { e.Network = value }
}

// OptEventAddr sets a field on the event.
func OptEventAddr(value string) EventOption {
	return func(e *Event) { e.Addr = value }
}

// OptEventDB sets a field on the event.
func OptEventDB(value string) EventOption {
	return func(e *Event) { e.DB = value }
}

// OptEventAuthUser sets a field on the event.
func OptEventAuthUser(value string) EventOption {
	return func(e *Event) { e.AuthUser = value }
}

// OptEventOp sets a field on the event.
func OptEventOp(value string) EventOption {
	return func(e *Event) { e.Op = value }
}

// OptEventArgs sets a field on the event.
func OptEventArgs(values ...string) EventOption {
	return func(e *Event) { e.Args = values }
}

// OptEventElapsed sets a field on the event.
func OptEventElapsed(value time.Duration) EventOption {
	return func(e *Event) { e.Elapsed = value }
}

// OptEventErr sets a field on the event.
func OptEventErr(value error) EventOption {
	return func(e *Event) { e.Err = value }
}

// Event represents a call to redis.
type Event struct {
	Network		string
	Addr		string
	AuthUser	string
	DB		string
	Op		string
	Args		[]string
	Elapsed		time.Duration
	Err		error
}

// GetFlag implements Event.
func (e Event) GetFlag() string	{ return Flag }

// WriteText writes the event text to the output.
func (e Event) WriteText(tf logger.TextFormatter, wr io.Writer) {
	fmt.Fprint(wr, "[")
	if len(e.AuthUser) > 0 {
		fmt.Fprint(wr, tf.Colorize(e.AuthUser, ansi.ColorLightWhite))
		fmt.Fprint(wr, "@")
	}
	if len(e.Addr) > 0 {
		fmt.Fprint(wr, tf.Colorize(e.Addr, ansi.ColorLightWhite))
	}
	if e.DB != "" {
		fmt.Fprint(wr, "/")
		fmt.Fprint(wr, tf.Colorize(fmt.Sprint(e.DB), ansi.ColorLightWhite))
	}
	fmt.Fprint(wr, "]")

	if len(e.Op) > 0 {
		fmt.Fprint(wr, logger.Space)
		fmt.Fprintf(wr, "[%s]", tf.Colorize(e.Op, ansi.ColorLightBlue))
	}
	if len(e.Args) > 0 {
		fmt.Fprint(wr, logger.Space)
		fmt.Fprintf(wr, "%s", strings.Join(e.Args, ", "))
	}

	fmt.Fprint(wr, logger.Space)
	fmt.Fprint(wr, e.Elapsed.String())

	if e.Err != nil {
		fmt.Fprint(wr, logger.Space)
		fmt.Fprint(wr, tf.Colorize("failed", ansi.ColorRed))
	}
}

// Decompose implements JSONWritable.
func (e Event) Decompose() map[string]interface{} {
	return map[string]interface{}{
		"addr":		e.Addr,
		"db":		e.DB,
		"op":		e.Op,
		"args":		e.Args,
		"elapsed":	timeutil.Milliseconds(e.Elapsed),
		"err":		e.Err,
	}
}
