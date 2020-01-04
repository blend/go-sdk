package r2

import (
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
	Flag = "http.client.request"
	// FlagResponse is a logger event flag.
	FlagResponse = "http.client.response"
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

// Event is a response to outgoing requests.
type Event struct {
	Flag string
	// The request metadata.
	Request *http.Request
	// The response metadata (excluding the body).
	Response *http.Response
	// The response body.
	Body []byte
	// Elapsed is the time elapsed.
	Elapsed time.Duration
}

// GetFlag implements logger.Event.
func (e Event) GetFlag() string { return e.Flag }

// WriteText writes the event to a text writer.
func (e Event) WriteText(tf logger.TextFormatter, wr io.Writer) {
	if e.Request != nil && e.Response != nil {
		io.WriteString(wr, fmt.Sprintf("%s %s %s (%v)", e.Request.Method, e.Request.URL.String(), webutil.ColorizeStatusCodeWithFormatter(tf, e.Response.StatusCode), e.Elapsed))
	} else if e.Request != nil {
		io.WriteString(wr, fmt.Sprintf("%s %s", e.Request.Method, e.Request.URL.String()))
	}
	if e.Body != nil {
		io.WriteString(wr, logger.Newline)
		io.WriteString(wr, string(e.Body))
	}
}

// Decompose implements logger.JSONWritable.
func (e *Event) Decompose() map[string]interface{} {
	output := make(map[string]interface{})
	if e.Request != nil {
		var url string
		if e.Request.URL != nil {
			url = e.Request.URL.String()
		}
		output["req"] = map[string]interface{}{
			"method":  e.Request.Method,
			"url":     url,
			"headers": e.Request.Header,
		}
	}
	if e.Response != nil {
		output["res"] = map[string]interface{}{
			"statusCode":      e.Response.StatusCode,
			"contentLength":   e.Response.ContentLength,
			"contentType":     tryHeader(e.Response.Header, "Content-Type", "content-type"),
			"contentEncoding": tryHeader(e.Response.Header, "Content-Encoding", "content-encoding"),
			"headers":         e.Response.Header,
			"cert":            webutil.ParseCertInfo(e.Response),
			"elapsed":         timeutil.Milliseconds(e.Elapsed),
		}
	}
	if e.Body != nil {
		output["body"] = string(e.Body)
	}

	return output
}

// EventJSONSchema is the json schema of the logger event.
type EventJSONSchema struct {
	Req struct {
		StartTime time.Time           `json:"startTime"`
		Method    string              `json:"method"`
		URL       string              `json:"url"`
		Headers   map[string][]string `json:"headers"`
	} `json:"req"`
	Res struct {
		CompleteTime  time.Time           `json:"completeTime"`
		StatusCode    int                 `json:"statusCode"`
		ContentLength int                 `json:"contentLength"`
		Headers       map[string][]string `json:"headers"`
	} `json:"res"`
	Body string `json:"body"`
}

func tryHeader(headers http.Header, keys ...string) string {
	for _, key := range keys {
		if values, hasValues := headers[key]; hasValues {
			return strings.Join(values, ";")
		}
	}
	return ""
}
