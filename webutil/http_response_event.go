package webutil

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/stringutil"
	"github.com/blend/go-sdk/timeutil"
)

var (
	_ logger.Event        = (*HTTPResponseEvent)(nil)
	_ logger.TextWritable = (*HTTPResponseEvent)(nil)
	_ logger.JSONWritable = (*HTTPResponseEvent)(nil)
)

// NewHTTPResponseEvent is an event representing a response to an http request.
func NewHTTPResponseEvent(req *http.Request, options ...HTTPResponseEventOption) HTTPResponseEvent {
	hre := HTTPResponseEvent{
		Request: req,
	}
	for _, option := range options {
		option(&hre)
	}
	return hre
}

// NewHTTPResponseEventListener returns a new web request event listener.
func NewHTTPResponseEventListener(listener func(context.Context, HTTPResponseEvent)) logger.Listener {
	return func(ctx context.Context, e logger.Event) {
		if typed, isTyped := e.(HTTPResponseEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// HTTPResponseEventOption is a function that modifies an http response event.
type HTTPResponseEventOption func(*HTTPResponseEvent)

// OptHTTPResponseRequest sets a field.
func OptHTTPResponseRequest(req *http.Request) HTTPResponseEventOption {
	return func(hre *HTTPResponseEvent) { hre.Request = req }
}

// OptHTTPResponseRoute sets a field.
func OptHTTPResponseRoute(route string) HTTPResponseEventOption {
	return func(hre *HTTPResponseEvent) { hre.Route = route }
}

// OptHTTPResponseContentLength sets a field.
func OptHTTPResponseContentLength(contentLength int) HTTPResponseEventOption {
	return func(hre *HTTPResponseEvent) { hre.ContentLength = contentLength }
}

// OptHTTPResponseContentType sets a field.
func OptHTTPResponseContentType(contentType string) HTTPResponseEventOption {
	return func(hre *HTTPResponseEvent) { hre.ContentType = contentType }
}

// OptHTTPResponseContentEncoding sets a field.
func OptHTTPResponseContentEncoding(contentEncoding string) HTTPResponseEventOption {
	return func(hre *HTTPResponseEvent) { hre.ContentEncoding = contentEncoding }
}

// OptHTTPResponseStatusCode sets a field.
func OptHTTPResponseStatusCode(statusCode int) HTTPResponseEventOption {
	return func(hre *HTTPResponseEvent) { hre.StatusCode = statusCode }
}

// OptHTTPResponseElapsed sets a field.
func OptHTTPResponseElapsed(elapsed time.Duration) HTTPResponseEventOption {
	return func(hre *HTTPResponseEvent) { hre.Elapsed = elapsed }
}

// OptHTTPResponseHeader sets a field.
func OptHTTPResponseHeader(header http.Header) HTTPResponseEventOption {
	return func(hre *HTTPResponseEvent) { hre.Header = header }
}

// HTTPResponseEvent is an event type for responses.
type HTTPResponseEvent struct {
	Request         *http.Request
	Route           string
	ContentLength   int
	ContentType     string
	ContentEncoding string
	StatusCode      int
	Elapsed         time.Duration
	Header          http.Header
}

// GetFlag implements event.
func (e HTTPResponseEvent) GetFlag() string { return HTTPResponse }

// WriteText implements TextWritable.
func (e HTTPResponseEvent) WriteText(tf logger.TextFormatter, wr io.Writer) {
	if ip := GetRemoteAddr(e.Request); len(ip) > 0 {
		io.WriteString(wr, ip)
		io.WriteString(wr, logger.Space)
	}
	io.WriteString(wr, tf.Colorize(e.Request.Method, ansi.ColorBlue))
	io.WriteString(wr, logger.Space)
	io.WriteString(wr, e.Request.URL.String())
	io.WriteString(wr, logger.Space)
	io.WriteString(wr, ColorizeStatusCodeWithFormatter(tf, e.StatusCode))
	io.WriteString(wr, logger.Space)
	io.WriteString(wr, e.Elapsed.String())
	if len(e.ContentType) > 0 {
		io.WriteString(wr, logger.Space)
		io.WriteString(wr, e.ContentType)
	}
	io.WriteString(wr, logger.Space)
	io.WriteString(wr, stringutil.FileSize(e.ContentLength))
}

// Decompose implements JSONWritable.
func (e HTTPResponseEvent) Decompose() map[string]interface{} {
	return map[string]interface{}{
		"ip":              GetRemoteAddr(e.Request),
		"userAgent":       GetUserAgent(e.Request),
		"verb":            e.Request.Method,
		"path":            e.Request.URL.Path,
		"route":           e.Route,
		"query":           e.Request.URL.RawQuery,
		"host":            e.Request.Host,
		"contentLength":   e.ContentLength,
		"contentType":     e.ContentType,
		"contentEncoding": e.ContentEncoding,
		"statusCode":      e.StatusCode,
		"elapsed":         timeutil.Milliseconds(e.Elapsed),
	}
}
