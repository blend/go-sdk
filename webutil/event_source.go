package webutil

import (
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/blend/go-sdk/stringutil"

	"github.com/blend/go-sdk/ex"
)

// NewEventSource returns a new event source.
// It is critical the response is *NOT* gzipped, as this will prevent `EventSource`
// from being able to effectively flush events.
func NewEventSource(output http.ResponseWriter) *EventSource {
	if _, ok := output.(*GZipResponseWriter); ok {
		panic("cannot create EventSource for GZipResponseWriter")
	}
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
	es.output.Header().Set(HeaderCacheControl, "no-cache")
	es.output.WriteHeader(http.StatusOK)

	err := es.eventUnsafe("ping")
	if err != nil {
		return err
	}
	return es.finishEventUnsafe()
}

// Ping sends the ping heartbeat event.
func (es *EventSource) Ping() error {
	return es.Event("ping")
}

// Event writes an event.
func (es *EventSource) Event(name string) error {
	es.Lock()
	defer es.Unlock()

	err := es.eventUnsafe(name)
	if err != nil {
		return err
	}
	return es.finishEventUnsafe()
}

// Data writes a data segment.
// It will slit lines on newline across multiple data events.
func (es *EventSource) Data(data string) error {
	es.Lock()
	defer es.Unlock()
	err := es.dataUnsafe(data)
	if err != nil {
		return err
	}
	return es.finishEventUnsafe()
}

// EventData sends an event with a given set of data.
func (es *EventSource) EventData(name, data string) error {
	es.Lock()
	defer es.Unlock()

	err := es.eventUnsafe(name)
	if err != nil {
		return err
	}
	err = es.dataUnsafe(data)
	if err != nil {
		return err
	}
	return es.finishEventUnsafe()
}

// EventDataWithID sends a named event with a given set of data and message identifier.
func (es *EventSource) EventDataWithID(name, data, id string) error {
	es.Lock()
	defer es.Unlock()

	err := es.eventUnsafe(name)
	if err != nil {
		return err
	}
	err = es.dataUnsafe(data)
	if err != nil {
		return err
	}
	err = es.idUnsafe(id)
	if err != nil {
		return err
	}
	return es.finishEventUnsafe()
}

//
// unsafe methods
//

func (es *EventSource) eventUnsafe(name string) error {
	_, err := io.WriteString(es.output, "event: "+strings.TrimSpace(name)+"\n")
	if err != nil {
		return ex.New(err)
	}
	return nil
}

// dataUnsafe will send one or many data message lines.
// if the `data` parameter contains newlines, the contents will be split up across multiple
// `data:` events to the client.
func (es *EventSource) dataUnsafe(data string) error {
	lines := stringutil.SplitLines(data, stringutil.OptSplitLinesIncludeEmptyLines(true))
	var err error
	for _, line := range lines {
		_, err = io.WriteString(es.output, "data: "+line+"\n")
		if err != nil {
			return ex.New(err)
		}
	}
	return nil
}

func (es *EventSource) idUnsafe(id string) error {
	_, err := io.WriteString(es.output, "id: "+strings.TrimSpace(id)+"\n")
	if err != nil {
		return ex.New(err)
	}
	return nil
}

// finishEventUnsafe writes a final `\n` or newline, and flushes the underlying http response.
func (es *EventSource) finishEventUnsafe() error {
	_, err := io.WriteString(es.output, "\n")
	if err != nil {
		return ex.New(err)
	}
	if typed, ok := es.output.(http.Flusher); ok {
		typed.Flush()
	}
	return nil
}
