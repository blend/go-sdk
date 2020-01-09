package webutil

import (
	"io"
	"net/http"
	"sync"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/stringutil"
)

// NewEventSource returns a new event source.
func NewEventSource(output http.ResponseWriter) *EventSource {
	return &EventSource{output: output}
}

// EventSource is a helper for writing event source info.
type EventSource struct {
	sync.Mutex
	output http.ResponseWriter
}

// StartSession starts an event source session.
func (es *EventSource) StartSession() error {
	es.Lock()
	defer es.Unlock()

	es.output.Header().Set(HeaderContentType, "text/event-stream")
	es.output.Header().Set(HeaderVary, "Content-Type")
	es.output.WriteHeader(http.StatusOK)
	return es.eventUnsafe("ping")
}

// Ping sends the ping heartbeat event.
func (es *EventSource) Ping() error {
	return es.Event("ping")
}

// Event writes an event.
func (es *EventSource) Event(name string) error {
	es.Lock()
	defer es.Unlock()
	return es.eventUnsafe(name)
}

// Data writes a data segment.
// It will slit lines on newline across multiple data events.
func (es *EventSource) Data(data string) error {
	es.Lock()
	defer es.Unlock()
	return es.dataUnsafe(data)
}

// EventData sends an event with a given set of data.
func (es *EventSource) EventData(name, data string) error {
	es.Lock()
	defer es.Unlock()
	_, err := io.WriteString(es.output, "event: "+name+"\n")
	if err != nil {
		return ex.New(err)
	}
	return es.dataUnsafe(data)
}

//
// unsafe methods
//

func (es *EventSource) eventUnsafe(name string) error {
	_, err := io.WriteString(es.output, "event: "+name+"\n\n")
	if err != nil {
		return ex.New(err)
	}
	if typed, ok := es.output.(http.Flusher); ok {
		typed.Flush()
	}
	return nil
}

func (es *EventSource) dataUnsafe(data string) error {
	lines := stringutil.SplitLines(data, stringutil.OptSplitLinesIncludeEmptyLines(true))
	for _, line := range lines {
		_, err := io.WriteString(es.output, "data: "+line+"\n")
		if err != nil {
			return ex.New(err)
		}
	}
	_, err := io.WriteString(es.output, "\n")
	if err != nil {
		return ex.New(err)
	}
	if typed, ok := es.output.(http.Flusher); ok {
		typed.Flush()
	}
	return nil
}
