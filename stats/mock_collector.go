package stats

import (
	"fmt"
	"time"

	"github.com/blend/go-sdk/timeutil"
)

// Assert that the mock collector implements Collector.
var (
	_ Collector = (*MockCollector)(nil)
)

// NewMockCollector returns a new mock collector.
func NewMockCollector() *MockCollector {
	return &MockCollector{
		Events: make(chan MockMetric, 32),
		Errors: make(chan error, 32),
	}
}

// MockCollector is a mocked collector for stats.
type MockCollector struct {
	namespace   string
	defaultTags []string

	Events chan MockMetric
	Errors chan error
}

// AddDefaultTag adds a default tag.
func (mc *MockCollector) AddDefaultTag(key, value string) {
	mc.defaultTags = append(mc.defaultTags, fmt.Sprintf("%s:%s", key, value))
}

// DefaultTags returns the default tags set.
func (mc MockCollector) DefaultTags() []string {
	return mc.defaultTags
}

// Count adds a mock count event to the event stream.
func (mc MockCollector) Count(name string, value int64, tags ...string) error {
	mc.Events <- MockMetric{Name: name, Count: value, Tags: append(mc.defaultTags, tags...)}
	if len(mc.Errors) > 0 {
		return <-mc.Errors
	}
	return nil
}

// Increment adds a mock count event to the event stream with value (1).
func (mc MockCollector) Increment(name string, tags ...string) error {
	mc.Events <- MockMetric{Name: name, Count: 1, Tags: append(mc.defaultTags, tags...)}
	if len(mc.Errors) > 0 {
		return <-mc.Errors
	}
	return nil
}

// Gauge adds a mock count event to the event stream with value (1).
func (mc MockCollector) Gauge(name string, value float64, tags ...string) error {
	mc.Events <- MockMetric{Name: name, Gauge: value, Tags: append(mc.defaultTags, tags...)}
	if len(mc.Errors) > 0 {
		return <-mc.Errors
	}
	return nil
}

// Histogram adds a mock count event to the event stream with value (1).
func (mc MockCollector) Histogram(name string, value float64, tags ...string) error {
	mc.Events <- MockMetric{Name: name, Histogram: value, Tags: append(mc.defaultTags, tags...)}
	if len(mc.Errors) > 0 {
		return <-mc.Errors
	}
	return nil
}

// TimeInMilliseconds adds a mock time in millis event to the event stream with a value.
func (mc MockCollector) TimeInMilliseconds(name string, value time.Duration, tags ...string) error {
	mc.Events <- MockMetric{Name: name, TimeInMilliseconds: timeutil.Milliseconds(value), Tags: append(mc.defaultTags, tags...)}
	if len(mc.Errors) > 0 {
		return <-mc.Errors
	}
	return nil
}

// MockMetric is a mock metric.
type MockMetric struct {
	Name               string
	Count              int64
	Gauge              float64
	Histogram          float64
	TimeInMilliseconds float64
	Tags               []string
}
