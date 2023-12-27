/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package r2

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/timeutil"
	"github.com/blend/go-sdk/webutil"
)

const (
	// Flag is a logger event flag.
	Flag	= "http.client.request"
	// FlagResponse is a logger event flag.
	FlagResponse	= "http.client.response"
)

// NewEvent returns a new event.
func NewEvent(flag string, options ...EventOption) Event {
	e := Event{
		Flag: flag,
	}
	for _, option := range options {
		option(&e)
	}
	return e
}

// NewEventListener returns a new r2 event listener.
func NewEventListener(listener func(context.Context, Event)) logger.Listener {
	return func(ctx context.Context, e logger.Event) {
		if typed, isTyped := e.(Event); isTyped {
			listener(ctx, typed)
		}
	}
}

// NewEventFilter returns a new r2 event filter.
func NewEventFilter(filter func(context.Context, Event) (Event, bool)) logger.Filter {
	return func(ctx context.Context, e logger.Event) (logger.Event, bool) {
		if typed, isTyped := e.(Event); isTyped {
			return filter(ctx, typed)
		}
		return e, false
	}
}

var (
	_	logger.Event		= (*Event)(nil)
	_	logger.TextWritable	= (*Event)(nil)
	_	logger.JSONWritable	= (*Event)(nil)
)

// Event is a response to outgoing requests.
type Event struct {
	Flag	string
	// The request metadata.
	Request	*http.Request
	// The response metadata (excluding the body).
	Response	*http.Response
	// The response body.
	Body	[]byte
	// Elapsed is the time elapsed.
	Elapsed	time.Duration
}

// GetFlag implements logger.Event.
func (e Event) GetFlag() string	{ return e.Flag }

// WriteText writes the event to a text writer.
func (e Event) WriteText(tf logger.TextFormatter, wr io.Writer) {
	if e.Request != nil && e.Response != nil {
		fmt.Fprintf(wr, "%s %s %s (%v)", e.Request.Method, e.Request.URL.String(), webutil.ColorizeStatusCodeWithFormatter(tf, e.Response.StatusCode), e.Elapsed)
	} else if e.Request != nil {
		fmt.Fprintf(wr, "%s %s", e.Request.Method, e.Request.URL.String())
	}
	if e.Body != nil {
		fmt.Fprint(wr, logger.Newline)
		fmt.Fprint(wr, string(e.Body))
	}
}

// Decompose implements logger.JSONWritable.
func (e Event) Decompose() map[string]interface{} {
	output := make(map[string]interface{})
	if e.Request != nil {
		var url string
		if e.Request.URL != nil {
			url = e.Request.URL.String()
		}
		output["req"] = map[string]interface{}{
			"method":	e.Request.Method,
			"url":		url,
		}
	}
	if e.Response != nil {
		output["res"] = map[string]interface{}{
			"statusCode":		e.Response.StatusCode,
			"contentLength":	e.Response.ContentLength,
			"contentType":		tryHeader(e.Response.Header, "Content-Type", "content-type"),
			"contentEncoding":	tryHeader(e.Response.Header, "Content-Encoding", "content-encoding"),
			"cert":			webutil.ParseCertInfo(e.Response),
			"elapsed":		timeutil.Milliseconds(e.Elapsed),
		}
	}
	if e.Body != nil {
		output["body"] = string(e.Body)
	}

	return output
}

// EventJSONSchema is the json schema of the logger event.
type EventJSONSchema struct {
	Req	struct {
		StartTime	time.Time		`json:"startTime"`
		Method		string			`json:"method"`
		URL		string			`json:"url"`
		Headers		map[string][]string	`json:"headers"`
	}	`json:"req"`
	Res	struct {
		CompleteTime	time.Time		`json:"completeTime"`
		StatusCode	int			`json:"statusCode"`
		ContentLength	int			`json:"contentLength"`
		Headers		map[string][]string	`json:"headers"`
	}	`json:"res"`
	Body	string	`json:"body"`
}

func tryHeader(headers http.Header, keys ...string) string {
	for _, key := range keys {
		if values, hasValues := headers[key]; hasValues {
			return strings.Join(values, ";")
		}
	}
	return ""
}
