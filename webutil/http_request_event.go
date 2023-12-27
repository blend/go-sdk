/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stringutil"
	"github.com/blend/go-sdk/timeutil"
)

var (
	_	logger.Event		= (*HTTPRequestEvent)(nil)
	_	logger.TextWritable	= (*HTTPRequestEvent)(nil)
	_	logger.JSONWritable	= (*HTTPRequestEvent)(nil)
)

// NewHTTPRequestEvent is an event representing a request to an http server.
func NewHTTPRequestEvent(req *http.Request, options ...HTTPRequestEventOption) HTTPRequestEvent {
	hre := HTTPRequestEvent{
		Request: req,
	}
	for _, option := range options {
		option(&hre)
	}
	return hre
}

// NewHTTPRequestEventListener returns a new web request event listener.
func NewHTTPRequestEventListener(listener func(context.Context, HTTPRequestEvent)) logger.Listener {
	return func(ctx context.Context, e logger.Event) {
		if typed, isTyped := e.(HTTPRequestEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// NewHTTPRequestEventFilter returns a new http request event filter.
func NewHTTPRequestEventFilter(filter func(context.Context, HTTPRequestEvent) (HTTPRequestEvent, bool)) logger.Filter {
	return func(ctx context.Context, e logger.Event) (logger.Event, bool) {
		if typed, isTyped := e.(HTTPRequestEvent); isTyped {
			return filter(ctx, typed)
		}
		return e, false
	}
}

// HTTPRequestEventOption is a function that modifies an http request event.
type HTTPRequestEventOption func(*HTTPRequestEvent)

// OptHTTPRequestRequest sets a field.
func OptHTTPRequestRequest(req *http.Request) HTTPRequestEventOption {
	return func(hre *HTTPRequestEvent) { hre.Request = req }
}

// OptHTTPRequestRoute sets a field.
func OptHTTPRequestRoute(route string) HTTPRequestEventOption {
	return func(hre *HTTPRequestEvent) { hre.Route = route }
}

// OptHTTPRequestContentLength sets a field.
func OptHTTPRequestContentLength(contentLength int) HTTPRequestEventOption {
	return func(hre *HTTPRequestEvent) { hre.ContentLength = contentLength }
}

// OptHTTPRequestContentType sets a field.
func OptHTTPRequestContentType(contentType string) HTTPRequestEventOption {
	return func(hre *HTTPRequestEvent) { hre.ContentType = contentType }
}

// OptHTTPRequestContentEncoding sets a field.
func OptHTTPRequestContentEncoding(contentEncoding string) HTTPRequestEventOption {
	return func(hre *HTTPRequestEvent) { hre.ContentEncoding = contentEncoding }
}

// OptHTTPRequestStatusCode sets a field.
func OptHTTPRequestStatusCode(statusCode int) HTTPRequestEventOption {
	return func(hre *HTTPRequestEvent) { hre.StatusCode = statusCode }
}

// OptHTTPRequestElapsed sets a field.
func OptHTTPRequestElapsed(elapsed time.Duration) HTTPRequestEventOption {
	return func(hre *HTTPRequestEvent) { hre.Elapsed = elapsed }
}

// OptHTTPRequestHeader sets a field.
func OptHTTPRequestHeader(header http.Header) HTTPRequestEventOption {
	return func(hre *HTTPRequestEvent) { hre.Header = header }
}

// HTTPRequestEvent is an event type for http requests.
type HTTPRequestEvent struct {
	Request		*http.Request
	Route		string
	ContentLength	int
	ContentType	string
	ContentEncoding	string
	StatusCode	int
	Elapsed		time.Duration
	Header		http.Header
}

// GetFlag implements event.
func (e HTTPRequestEvent) GetFlag() string	{ return FlagHTTPRequest }

// WriteText implements TextWritable.
func (e HTTPRequestEvent) WriteText(tf logger.TextFormatter, wr io.Writer) {
	if ip := GetRemoteAddr(e.Request); len(ip) > 0 {
		fmt.Fprint(wr, ip)
		fmt.Fprint(wr, logger.Space)
	}
	fmt.Fprint(wr, tf.Colorize(e.Request.Method, ansi.ColorBlue))
	fmt.Fprint(wr, logger.Space)
	fmt.Fprint(wr, e.Request.URL.String())
	fmt.Fprint(wr, logger.Space)
	fmt.Fprint(wr, ColorizeStatusCodeWithFormatter(tf, e.StatusCode))
	fmt.Fprint(wr, logger.Space)
	fmt.Fprint(wr, e.Elapsed.String())
	if len(e.ContentType) > 0 {
		fmt.Fprint(wr, logger.Space)
		fmt.Fprint(wr, e.ContentType)
	}
	fmt.Fprint(wr, logger.Space)
	fmt.Fprint(wr, stringutil.FileSize(e.ContentLength))
}

// Decompose implements JSONWritable.
func (e HTTPRequestEvent) Decompose() map[string]interface{} {
	return map[string]interface{}{
		"ip":			GetRemoteAddr(e.Request),
		"userAgent":		GetUserAgent(e.Request),
		"verb":			e.Request.Method,
		"path":			e.Request.URL.Path,
		"route":		e.Route,
		"query":		e.Request.URL.RawQuery,
		"host":			e.Request.Host,
		"contentLength":	e.ContentLength,
		"contentType":		e.ContentType,
		"contentEncoding":	e.ContentEncoding,
		"statusCode":		e.StatusCode,
		"elapsed":		timeutil.Milliseconds(e.Elapsed),
	}
}
