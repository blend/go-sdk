package logger

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/blend/go-sdk/webutil"
)

// these are compile time assertions
var (
	_ Event          = (*HTTPRequestEvent)(nil)
	_ TextWritable   = (*HTTPRequestEvent)(nil)
	_ json.Marshaler = (*HTTPRequestEvent)(nil)
)

// NewHTTPRequestEvent creates a new web request event.
func NewHTTPRequestEvent(req *http.Request, options ...HTTPRequestEventOption) *HTTPRequestEvent {
	hre := &HTTPRequestEvent{
		EventMeta: NewEventMeta(HTTPRequest),
		Request:   req,
	}
	for _, option := range options {
		option(hre)
	}
	return hre
}

// NewHTTPRequestEventListener returns a new web request event listener.
func NewHTTPRequestEventListener(listener func(context.Context, *HTTPRequestEvent)) Listener {
	return func(ctx context.Context, e Event) {
		if typed, isTyped := e.(*HTTPRequestEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// HTTPRequestEventOption sets a field on an HTTPRequestEventOption.
type HTTPRequestEventOption func(*HTTPRequestEvent)

// OptHTTPRequestMeta sets a field on an HTTPRequestEvent.
func OptHTTPRequestMeta(options ...EventMetaOption) HTTPRequestEventOption {
	return func(hre *HTTPRequestEvent) {
		for _, option := range options {
			option(hre.EventMeta)
		}
	}
}

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

// OptHTTPRequestState sets a field on an HTTPRequestEvent.
func OptHTTPRequestState(state map[interface{}]interface{}) HTTPRequestEventOption {
	return func(hre *HTTPRequestEvent) {
		hre.State = state
	}
}

// HTTPRequestEvent is an event type for http responses.
type HTTPRequestEvent struct {
	*EventMeta
	Request *http.Request
	Route   string
	State   map[interface{}]interface{}
}

// WriteText implements TextWritable.
func (e *HTTPRequestEvent) WriteText(formatter TextFormatter, wr io.Writer) {
	WriteHTTPRequest(formatter, wr, e.Request)
}

// MarshalJSON marshals the event as json.
func (e *HTTPRequestEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(MergeDecomposed(e.EventMeta.Decompose(), map[string]interface{}{
		"verb":      e.Request.Method,
		"path":      e.Request.URL.Path,
		"host":      e.Request.Host,
		"ip":        webutil.GetRemoteAddr(e.Request),
		"userAgent": webutil.GetUserAgent(e.Request),
	}))
}
