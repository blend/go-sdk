package webutil

import (
	"io"
	"net/http"

	"github.com/blend/go-sdk/ex"
)

// EventSource is a helper for writing event source info.
type EventSource struct {
	Output http.ResponseWriter
}

// StartSession starts an event source session.
func (es EventSource) StartSession() error {
	es.Output.Header().Set(HeaderContentType, "text/event-stream")
	es.Output.Header().Set(HeaderVary, "Content-Type")
	es.Output.WriteHeader(http.StatusOK)
	return es.Ping()
}

// Ping sends the ping heartbeat event.
func (es EventSource) Ping() error {
	return es.Event("ping")
}

// Event writes an event.
func (es EventSource) Event(name string) error {
	_, err := io.WriteString(es.Output, "event: "+name+"\n\n")
	if err != nil {
		return ex.New(err)
	}
	if typed, ok := es.Output.(http.Flusher); ok {
		typed.Flush()
	}
	return nil
}

// Data writes a data event.
func (es EventSource) Data(data string) error {
	_, err := io.WriteString(es.Output, "data: "+data+"\n\n")
	if err != nil {
		return ex.New(err)
	}
	if typed, ok := es.Output.(http.Flusher); ok {
		typed.Flush()
	}
	return nil
}

// EventData sends an event with a given set of data.
func (es EventSource) EventData(name, data string) error {
	_, err := io.WriteString(es.Output, "event: "+name+"\n")
	if err != nil {
		return ex.New(err)
	}
	_, err = io.WriteString(es.Output, "data: "+data+"\n\n")
	if err != nil {
		return ex.New(err)
	}
	if typed, ok := es.Output.(http.Flusher); ok {
		typed.Flush()
	}
	return nil
}
