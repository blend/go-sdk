package logger

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/blend/go-sdk/timeutil"
	"github.com/blend/go-sdk/webutil"
)

var (
	_ Event          = (*HTTPResponseEvent)(nil)
	_ TextWritable   = (*HTTPResponseEvent)(nil)
	_ json.Marshaler = (*HTTPResponseEvent)(nil)
)

// NewHTTPResponseEvent is an event representing a response to an http request.
func NewHTTPResponseEvent(req *http.Request, options ...HTTPResponseEventOption) *HTTPResponseEvent {
	hre := &HTTPResponseEvent{
		EventMeta: NewEventMeta(HTTPResponse),
		Request:   req,
	}
	for _, option := range options {
		option(hre)
	}
	return hre
}

// NewHTTPResponseEventListener returns a new web request event listener.
func NewHTTPResponseEventListener(listener func(context.Context, *HTTPResponseEvent)) Listener {
	return func(ctx context.Context, e Event) {
		if typed, isTyped := e.(*HTTPResponseEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// HTTPResponseEventOption is a function that modifies an http response event.
type HTTPResponseEventOption func(*HTTPResponseEvent)

// OptHTTPResponseMeta sets a fields on the event meta.
func OptHTTPResponseMeta(options ...EventMetaOption) HTTPResponseEventOption {
	return func(hre *HTTPResponseEvent) {
		for _, option := range options {
			option(hre.EventMeta)
		}
	}
}

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

// OptHTTPResponseState sets a field.
func OptHTTPResponseState(state interface{}) HTTPResponseEventOption {
	return func(hre *HTTPResponseEvent) { hre.State = state }
}

// HTTPResponseEvent is an event type for responses.
type HTTPResponseEvent struct {
	*EventMeta
	Request         *http.Request
	Route           string
	ContentLength   int
	ContentType     string
	ContentEncoding string
	StatusCode      int
	Elapsed         time.Duration
	Header          http.Header
	State           interface{}
}

// WriteText implements TextWritable.
func (e HTTPResponseEvent) WriteText(formatter TextFormatter, wr io.Writer) {
	WriteHTTPResponse(formatter, wr, e.Request, e.StatusCode, e.ContentLength, e.ContentType, e.Elapsed)
}

// MarshalJSON implements json.Marshaler.
func (e HTTPResponseEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(MergeDecomposed(e.EventMeta.Decompose(), map[string]interface{}{
		"ip":              webutil.GetRemoteAddr(e.Request),
		"userAgent":       webutil.GetUserAgent(e.Request),
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
		"state":           e.State,
	}))
}
