/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stats

// Assert that the mock collector implements Collector.
var (
	_ EventCollector = (*MockEventCollector)(nil)
)

// NewMockEventCollector returns a new mock collector.
func NewMockEventCollector() *MockEventCollector {
	return &MockEventCollector{
		Events: make(chan Event),
	}
}

// MockEventCollector is a mocked collector for stats.
type MockEventCollector struct {
	defaultTags []string
	Events      chan Event
}

// AddDefaultTag adds a default tag.
func (mec *MockEventCollector) AddDefaultTag(name, value string) {
	mec.defaultTags = append(mec.defaultTags, Tag(name, value))
}

// AddDefaultTags adds default tags.
func (mec *MockEventCollector) AddDefaultTags(tags ...string) {
	mec.defaultTags = append(mec.defaultTags, tags...)
}

// DefaultTags returns the default tags set.
func (mec MockEventCollector) DefaultTags() []string {
	return mec.defaultTags
}

// SendEvent sends an event.
func (mec MockEventCollector) SendEvent(e Event) error {
	mec.Events <- e
	return nil
}

// CreateEvent creates a mock event with the default tags.
func (mec MockEventCollector) CreateEvent(title, text string, tags ...string) Event {
	return Event{
		Title: title,
		Text:  text,
		Tags:  append(mec.defaultTags, tags...),
	}
}
