package stats

import (
	"time"

	"github.com/blend/go-sdk/ex"
)

// MultiCollector is a class that wraps a set of statsd collectors
type MultiCollector []Collector

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
