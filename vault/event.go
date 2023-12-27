/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package vault

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/timeutil"
)

var (
	_	logger.Event		= (*Event)(nil)
	_	logger.TextWritable	= (*Event)(nil)
	_	logger.JSONWritable	= (*Event)(nil)
)

const (
	// Flag is the logger flag.
	Flag = "vault"
)

// NewEventListener returns a new logger listener for a given event.
func NewEventListener(action func(context.Context, Event)) logger.Listener {
	return func(ctx context.Context, e logger.Event) {
		if typed, ok := e.(Event); ok {
			action(ctx, typed)
		}
	}
}

// NewEvent returns a new event from a request.
func NewEvent(req *http.Request) *Event {
	return &Event{
		Remote:	req.URL.Host,
		Method:	req.Method,
		Path:	strings.TrimPrefix(req.URL.Path, "/v1/"),
	}
}

// Event is an event.
type Event struct {
	Remote		string
	Method		string
	Path		string
	StatusCode	int
	Elapsed		time.Duration
}

// GetFlag implements logger.Event.
func (e Event) GetFlag() string	{ return Flag }

// WriteText writes text for the event.
func (e *Event) WriteText(tf logger.TextFormatter, wr io.Writer) {
	fmt.Fprint(wr, "["+tf.Colorize(e.Method, ansi.ColorBlue)+"]")
	fmt.Fprint(wr, logger.Space)
	fmt.Fprint(wr, e.Remote)
	fmt.Fprint(wr, logger.Space)
	fmt.Fprint(wr, e.Path)
	fmt.Fprint(wr, logger.Space)
	fmt.Fprint(wr, e.StatusCode)
	fmt.Fprint(wr, logger.Space)
	fmt.Fprint(wr, e.Elapsed)
}

// Decompose impements logger.JSONWritable.
func (e *Event) Decompose() map[string]interface{} {
	return map[string]interface{}{
		"remote":	e.Remote,
		"method":	e.Method,
		"path":		e.Path,
		"statusCode":	e.StatusCode,
		"elapsed":	timeutil.Milliseconds(e.Elapsed),
	}
}
