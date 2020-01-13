package stats

import (
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestMockEventCollectorDefaultTags(t *testing.T) {
	assert := assert.New(t)

	assertTags := func(actualTags []string) {
		assert.Len(actualTags, 3)
		assert.Equal("k1:v1", actualTags[0])
		assert.Equal("k2:v2", actualTags[1])
		assert.Equal("k3:v3", actualTags[2])
	}

	collector := NewMockEventCollector()
	collector.AddDefaultTag("k1", "v1")
	collector.AddDefaultTag("k2", "v2")

	tags := collector.DefaultTags()
	assert.Len(tags, 2)
	assert.Equal("k1:v1", tags[0])
	assert.Equal("k2:v2", tags[1])

	event := collector.CreateEvent("event", "text", "k3:v3")
	assertTags(event.Tags)

	go collector.SendEvent(event)
	receivedEvent := <-collector.Events
	assertTags(receivedEvent.Tags)
}

func TestMockEventCollectorCreateEvent(t *testing.T) {
	assert := assert.New(t)

	collector := NewMockEventCollector()
	event := collector.CreateEvent("event", "text", "k:v")
	assert.Equal("event", event.Title)
	assert.Equal("text", event.Text)
	assert.Len(event.Tags, 1)
	assert.Equal("k:v", event.Tags[0])
}

func TestMockEventCollectorSendEvent(t *testing.T) {
	assert := assert.New(t)

	collector := NewMockEventCollector()
	event := collector.CreateEvent("event", "text", "k:v")

	go collector.SendEvent(event)
	receivedEvent := <-collector.Events
	assert.Equal("event", receivedEvent.Title)
	assert.Equal("text", receivedEvent.Text)
}
