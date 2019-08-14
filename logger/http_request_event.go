package logger

import (
	"context"
	"io"
	"net/http"

	"github.com/blend/go-sdk/webutil"
)

// these are compile time assertions
var (
	_ Event        = (*HTTPRequestEvent)(nil)
	_ TextWritable = (*HTTPRequestEvent)(nil)
	_ JSONWritable = (*HTTPRequestEvent)(nil)
)

// NewHTTPRequestEvent creates a new web request event.
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
func NewHTTPRequestEventListener(listener func(context.Context, HTTPRequestEvent)) Listener {
	return func(ctx context.Context, e Event) {
		if typed, isTyped := e.(HTTPRequestEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// HTTPRequestEventOption sets a field on an HTTPRequestEventOption.
type HTTPRequestEventOption func(*HTTPRequestEvent)

// OptHTTPRequest sets a field on an HTTPRequestEvent.
func OptHTTPRequest(req *http.Request) HTTPRequestEventOption {
	return func(hre *HTTPRequestEvent) {
		hre.Request = req
	}
}

// OptHTTPRequestRoute sets a field on an HTTPRequestEvent.
func OptHTTPRequestRoute(route string) HTTPRequestEventOption {
	return func(hre *HTTPRequestEvent) {
		hre.Route = route
	}
}

// HTTPRequestEvent is an event type for http responses.
type HTTPRequestEvent struct {
	Request *http.Request
	Route   string
}

// GetFlag implements Event.
func (e HTTPRequestEvent) GetFlag() string { return HTTPRequest }

// WriteText implements TextWritable.
func (e HTTPRequestEvent) WriteText(formatter TextFormatter, wr io.Writer) {
	WriteHTTPRequest(formatter, wr, e.Request)
}

// Decompose implements JSONWritable.
func (e HTTPRequestEvent) Decompose() map[string]interface{} {
	return map[string]interface{}{
		"verb":      e.Request.Method,
		"path":      e.Request.URL.Path,
		"query":     e.Request.URL.RawQuery,
		"host":      e.Request.Host,
		"route":     e.Route,
		"ip":        webutil.GetRemoteAddr(e.Request),
		"userAgent": webutil.GetUserAgent(e.Request),
	}
}
