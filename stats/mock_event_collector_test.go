/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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
	collector.AddDefaultTags(Tag("k1", "v1"))
	collector.AddDefaultTags(Tag("k2", "v2"))

	tags := collector.DefaultTags()
	assert.Len(tags, 2)
	assert.Equal("k1:v1", tags[0])
	assert.Equal("k2:v2", tags[1])

	event := collector.CreateEvent("event", "text", "k3:v3")
	assertTags(event.Tags)

	go func() { _ = collector.SendEvent(event) }()
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

	go func() { _ = collector.SendEvent(event) }()
	receivedEvent := <-collector.Events
	assert.Equal("event", receivedEvent.Title)
	assert.Equal("text", receivedEvent.Text)
}
