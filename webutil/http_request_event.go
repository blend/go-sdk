package webutil

import (
	"context"
	"io"
	"net/http"
	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/logger"
)

// these are compile time assertions
var (
	_ logger.Event        = (*HTTPRequestEvent)(nil)
	_ logger.TextWritable = (*HTTPRequestEvent)(nil)
	_ logger.JSONWritable = (*HTTPRequestEvent)(nil)
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
func NewHTTPRequestEventListener(listener func(context.Context, HTTPRequestEvent)) logger.Listener {
	return func(ctx context.Context, e logger.Event) {
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
func (e HTTPRequestEvent) WriteText(tf logger.TextFormatter, wr io.Writer) {
	if ip := GetRemoteAddr(e.Request); len(ip) > 0 {
		io.WriteString(wr, ip)
		io.WriteString(wr, logger.Space)
	}
	io.WriteString(wr, tf.Colorize(e.Request.Method, ansi.ColorBlue))
	if e.Request.URL != nil {
		io.WriteString(wr, logger.Space)
		io.WriteString(wr, e.Request.URL.String())
	}
}

// Decompose implements JSONWritable.
func (e HTTPRequestEvent) Decompose() map[string]interface{} {
	return map[string]interface{}{
		"verb":      e.Request.Method,
		"path":      e.Request.URL.Path,
		"query":     e.Request.URL.RawQuery,
		"host":      e.Request.Host,
		"route":     e.Route,
		"ip":        GetRemoteAddr(e.Request),
		"userAgent": GetUserAgent(e.Request),
	}
}
