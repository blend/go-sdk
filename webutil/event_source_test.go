package webutil

import (
	"bytes"
	"net/http"
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestEventSourceStartSession(t *testing.T) {
	assert := assert.New(t)

	buffer := new(bytes.Buffer)
	rw := NewMockResponse(buffer)
	es := NewEventSource(rw)
	err := es.StartSession()
	assert.Nil(err)

	assert.Equal(http.StatusOK, rw.StatusCode())
	assert.Equal("text/event-stream", rw.Header().Get(HeaderContentType))
	assert.Equal("Content-Type", rw.Header().Get(HeaderVary))
	assert.Equal("event: ping\n\n", buffer.String())
}

func TestEventSourcePing(t *testing.T) {
	assert := assert.New(t)

	buffer := new(bytes.Buffer)
	rw := NewMockResponse(buffer)
	es := NewEventSource(rw)
	assert.Nil(es.Ping())
	assert.Equal("event: ping\n\n", buffer.String())
}

func TestEventSourceEvent(t *testing.T) {
	assert := assert.New(t)

	buffer := new(bytes.Buffer)
	rw := NewMockResponse(buffer)
	es := NewEventSource(rw)
	assert.Nil(es.Event("test event"))
	assert.Equal("event: test event\n\n", buffer.String())
}

func TestEventSourceData(t *testing.T) {
	assert := assert.New(t)

	buffer := new(bytes.Buffer)
	rw := NewMockResponse(buffer)
	es := NewEventSource(rw)
	assert.Nil(es.Data("test event data"))
	assert.Equal("data: test event data\n\n", buffer.String())
}

func TestEventSourceDataLines(t *testing.T) {
	assert := assert.New(t)

	buffer := new(bytes.Buffer)
	rw := NewMockResponse(buffer)
	es := NewEventSource(rw)
	assert.Nil(es.Data("test event data one\ntest event data two\n"))
	assert.Equal("data: test event data one\ndata: test event data two\n\n", buffer.String())
}

func TestEventSourceEventData(t *testing.T) {
	assert := assert.New(t)

	buffer := new(bytes.Buffer)
	rw := NewMockResponse(buffer)
	es := NewEventSource(rw)
	assert.Nil(es.EventData("test event", "test event data"))
	assert.Equal("event: test event\ndata: test event data\n\n", buffer.String())
}

func TestEventSourceEventDataLines(t *testing.T) {
	assert := assert.New(t)

	buffer := new(bytes.Buffer)
	rw := NewMockResponse(buffer)
	es := NewEventSource(rw)
	assert.Nil(es.EventData("test event", "test event data one\ntest event data two\n"))
	assert.Equal("event: test event\ndata: test event data one\ndata: test event data two\n\n", buffer.String())
}
