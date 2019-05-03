package r2

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

const (
	// Flag is a logger event flag.
	Flag = "http.client.request"
	// FlagResponse is a logger event flag.
	FlagResponse = "http.client.response"
)

// NewEvent returns a new event.
func NewEvent(flag string, options ...EventOption) *Event {
	e := &Event{
		EventMeta: logger.NewEventMeta(flag),
	}
	for _, option := range options {
		option(e)
	}
	return e
}

// Event is a response to outgoing requests.
type Event struct {
	*logger.EventMeta

	// Started is the time the request was started.
	// It is used for elapsed time calculations.
	Started time.Time
	// The request metadata.
	Request *http.Request
	// The response metadata (excluding the body).
	Response *http.Response
	// The response body.
	Body []byte
}

// WriteText writes the event to a text writer.
func (e *Event) WriteText(tf logger.TextFormatter, wr io.Writer) {
	if e.Request != nil && e.Response != nil {
		io.WriteString(wr, fmt.Sprintf("%s %s %s (%v)", e.Request.Method, e.Request.URL.String(), logger.ColorizeStatusCodeWithFormatter(tf, e.Response.StatusCode), e.GetTimestamp().Sub(e.Started)))
	} else if e.Request != nil {
		io.WriteString(wr, fmt.Sprintf("%s %s", e.Request.Method, e.Request.URL.String()))
	}
	if e.Body != nil {
		io.WriteString(wr, logger.Newline)
		io.WriteString(wr, string(e.Body))
	}
}

// MarshalJSON implements json.Marshaler.
func (e *Event) MarshalJSON() ([]byte, error) {
	output := make(map[string]interface{})
	if e.Request != nil {
		var url string
		if e.Request.URL != nil {
			url = e.Request.URL.String()
		}
		output["req"] = map[string]interface{}{
			"startTime": e.Started,
			"method":    e.Request.Method,
			"url":       url,
			"headers":   e.Request.Header,
		}
	}
	if e.Response != nil {
		output["res"] = map[string]interface{}{
			"completeTime":    e.GetTimestamp(),
			"statusCode":      e.Response.StatusCode,
			"contentLength":   e.Response.ContentLength,
			"contentType":     tryHeader(e.Response.Header, "Content-Type", "content-type"),
			"contentEncoding": tryHeader(e.Response.Header, "Content-Encoding", "content-encoding"),
			"headers":         e.Response.Header,
			"cert":            webutil.ParseCertInfo(e.Response),
		}
	}
	if e.Body != nil {
		output["body"] = string(e.Body)
	}

	return json.Marshal(logger.MergeDecomposed(e.EventMeta.Decompose(), output))
}

func tryHeader(headers http.Header, keys ...string) string {
	for _, key := range keys {
		if values, hasValues := headers[key]; hasValues {
			return strings.Join(values, ";")
		}
	}
	return ""
}
