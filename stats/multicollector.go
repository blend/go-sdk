package stats

import (
	"sort"
	"time"

	"github.com/blend/go-sdk/ex"
)

var (
	_ Collector = (*MultiCollector)(nil)
)

// MultiCollector is a class that wraps a set of statsd collectors
type MultiCollector []Collector

// AddDefaultTag implements Taggable.
func (collectors MultiCollector) AddDefaultTag(name, value string) {
	for _, collector := range collectors {
		collector.AddDefaultTag(name, value)
	}
}

// DefaultTags returns the unique default tags for the collector group.
func (collectors MultiCollector) DefaultTags() (output []string) {
	values := map[string]bool{}
	for _, collector := range collectors {
		for _, tag := range collector.DefaultTags() {
			values[tag] = true
		}
	}
	for key := range values {
		output = append(output, key)
	}
	sort.Strings(output)
	return
}

// Count increments a counter by a value and writes to the different hosts
func (collectors MultiCollector) Count(name string, value int64, tags ...string) error {
	for _, collector := range collectors {
		if err := collector.Count(name, value, tags...); err != nil {
			return ex.New(err)
		}
	}
	return nil
}

// Increment increments a counter by 1 and writes to the different hosts
func (collectors MultiCollector) Increment(name string, tags ...string) error {
	for _, collector := range collectors {
		if err := collector.Increment(name, tags...); err != nil {
			return ex.New(err)
		}
	}
	return nil
}

// Gauge sets a gauge value and writes to the different hosts
func (collectors MultiCollector) Gauge(name string, value float64, tags ...string) error {
	for _, collector := range collectors {
		if err := collector.Gauge(name, value, tags...); err != nil {
			return ex.New(err)
		}
	}
	return nil
}

// Histogram sets a guage value and writes to the different hosts
func (collectors MultiCollector) Histogram(name string, value float64, tags ...string) error {
	for _, collector := range collectors {
		if err := collector.Histogram(name, value, tags...); err != nil {
			return ex.New(err)
		}
	}
	return nil
}

// TimeInMilliseconds sets a timing value and writes to the different hosts
func (collectors MultiCollector) TimeInMilliseconds(name string, value time.Duration, tags ...string) error {
	for _, collector := range collectors {
		if err := collector.TimeInMilliseconds(name, value, tags...); err != nil {
			return ex.New(err)
		}
	}
	return nil
}
