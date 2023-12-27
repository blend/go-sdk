/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stats

import (
	"time"

	"github.com/blend/go-sdk/timeutil"
)

// Assert that the mock collector implements Collector.
var (
	_ Collector = (*MockCollector)(nil)
)

// NewMockCollector returns a new mock collector.
func NewMockCollector(capacity int) *MockCollector {
	return &MockCollector{
		Metrics:     make(chan MockMetric, capacity),
		Events:      make(chan Event, capacity),
		Errors:      make(chan error, capacity),
		FlushErrors: make(chan error, capacity),
		CloseErrors: make(chan error, capacity),
	}
}

// MockCollector is a mocked collector for stats.
type MockCollector struct {
	Field struct {
		Namespace   string
		DefaultTags []string
	}

	Metrics     chan MockMetric
	Events      chan Event
	Errors      chan error
	FlushErrors chan error
	CloseErrors chan error
}

// AllMetrics returns all the metrics from the collector.
func (mc *MockCollector) AllMetrics() (output []MockMetric) {
	metricCount := len(mc.Metrics)
	for x := 0; x < metricCount; x++ {
		output = append(output, <-mc.Metrics)
	}
	return
}

// GetCount returns the number of events logged for a given metric name.
func (mc *MockCollector) GetCount(metricName string) (count int) {
	var metric MockMetric
	metricCount := len(mc.Metrics)
	for x := 0; x < metricCount; x++ {
		metric = <-mc.Metrics
		if metric.Name == metricName {
			count++
		}
		mc.Metrics <- metric
	}
	return
}

// GetTagCount returns the number of events logged for a given metric name & tag name.
func (mc *MockCollector) GetTagCount(metricName string, tagName string) (count int) {
	metrics := mc.AllMetrics()
	metricCount := len(metrics)
	for x := 0; x < metricCount; x++ {
		metric := metrics[x]
		if metric.Name == metricName {
			for i := 0; i < len(metric.Tags); i++ {
				if tagName == metric.Tags[i] {
					count++
				}
			}
		}
	}
	return
}

func (mc *MockCollector) makeName(name string) string {
	if mc.Field.Namespace != "" {
		return mc.Field.Namespace + name
	}
	return name
}

// AddDefaultTag adds a default tag.
func (mc *MockCollector) AddDefaultTag(name, value string) {
	mc.Field.DefaultTags = append(mc.Field.DefaultTags, Tag(name, value))
}

// AddDefaultTags adds default tags.
func (mc *MockCollector) AddDefaultTags(tags ...string) {
	mc.Field.DefaultTags = append(mc.Field.DefaultTags, tags...)
}

// DefaultTags returns the default tags set.
func (mc MockCollector) DefaultTags() []string {
	return mc.Field.DefaultTags
}

// Count adds a mock count event to the event stream.
func (mc MockCollector) Count(name string, value int64, tags ...string) error {
	mc.Metrics <- MockMetric{Name: mc.makeName(name), Count: value, Tags: append(mc.Field.DefaultTags, tags...)}
	if len(mc.Errors) > 0 {
		return <-mc.Errors
	}
	return nil
}

// Increment adds a mock count event to the event stream with value (1).
func (mc MockCollector) Increment(name string, tags ...string) error {
	mc.Metrics <- MockMetric{Name: mc.makeName(name), Count: 1, Tags: append(mc.Field.DefaultTags, tags...)}
	if len(mc.Errors) > 0 {
		return <-mc.Errors
	}
	return nil
}

// Gauge adds a mock count event to the event stream with value (1).
func (mc MockCollector) Gauge(name string, value float64, tags ...string) error {
	mc.Metrics <- MockMetric{Name: mc.makeName(name), Gauge: value, Tags: append(mc.Field.DefaultTags, tags...)}
	if len(mc.Errors) > 0 {
		return <-mc.Errors
	}
	return nil
}

// Histogram adds a mock count event to the event stream with value (1).
func (mc MockCollector) Histogram(name string, value float64, tags ...string) error {
	mc.Metrics <- MockMetric{Name: mc.makeName(name), Histogram: value, Tags: append(mc.Field.DefaultTags, tags...)}
	if len(mc.Errors) > 0 {
		return <-mc.Errors
	}
	return nil
}

// Distribution adds a mock count event to the event stream with value (1).
func (mc MockCollector) Distribution(name string, value float64, tags ...string) error {
	mc.Metrics <- MockMetric{Name: mc.makeName(name), Distribution: value, Tags: append(mc.Field.DefaultTags, tags...)}
	if len(mc.Errors) > 0 {
		return <-mc.Errors
	}
	return nil
}

// TimeInMilliseconds adds a mock time in millis event to the event stream with a value.
func (mc MockCollector) TimeInMilliseconds(name string, value time.Duration, tags ...string) error {
	mc.Metrics <- MockMetric{Name: mc.makeName(name), TimeInMilliseconds: timeutil.Milliseconds(value), Tags: append(mc.Field.DefaultTags, tags...)}
	if len(mc.Errors) > 0 {
		return <-mc.Errors
	}
	return nil
}

// Flush does nothing on a MockCollector.
func (mc MockCollector) Flush() error {
	if len(mc.FlushErrors) > 0 {
		return <-mc.FlushErrors
	}
	return nil
}

// Close returns an error from the errors channel if any.
func (mc MockCollector) Close() error {
	if len(mc.CloseErrors) > 0 {
		return <-mc.CloseErrors
	}
	return nil
}

// SendEvent sends an event.
func (mc MockCollector) SendEvent(e Event) error {
	mc.Events <- e
	return nil
}

// CreateEvent creates a mock event with the default tags.
func (mc MockCollector) CreateEvent(title, text string, tags ...string) Event {
	return Event{
		Title: title,
		Text:  text,
		Tags:  append(mc.Field.DefaultTags, tags...),
	}
}

// MockMetric is a mock metric.
type MockMetric struct {
	Name               string
	Count              int64
	Gauge              float64
	Histogram          float64
	Distribution       float64
	TimeInMilliseconds float64
	Tags               []string
}
