package stats

import (
	"sort"
	"strings"
	"time"
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

// AddDefaultTags implements Taggable.
func (collectors MultiCollector) AddDefaultTags(tags ...string) {
	for _, collector := range collectors {
		collector.AddDefaultTags(tags...)
	}
}

// DefaultTags returns the unique default tags for the collectors.
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

// HasTagKey returns if the collector has a given tag key
// in *any* collector's default tags.
func (collectors MultiCollector) HasTagKey(tagKey string) bool {
	var key, value string
	for _, collector := range collectors {
		for _, tag := range collector.DefaultTags() {
			key, value = SplitTag(tag)
			if strings.EqualFold(tagKey, key) && value != "" {
				return true
			}
		}
	}
	return false
}

// Count increments a counter by a value and writes to the collectors.
func (collectors MultiCollector) Count(name string, value int64, tags ...string) (err error) {
	for _, collector := range collectors {
		err = collector.Count(name, value, tags...)
		if err != nil {
			return
		}
	}
	return
}

// Increment increments a counter by 1 and writes to the collectors.
func (collectors MultiCollector) Increment(name string, tags ...string) (err error) {
	for _, collector := range collectors {
		err = collector.Increment(name, tags...)
		if err != nil {
			return
		}
	}
	return
}

// Gauge sets a gauge value and writes to the collectors.
func (collectors MultiCollector) Gauge(name string, value float64, tags ...string) (err error) {
	for _, collector := range collectors {
		err = collector.Gauge(name, value, tags...)
		if err != nil {
			return
		}
	}
	return
}

// Histogram sets a histogram value and writes to the collectors.
func (collectors MultiCollector) Histogram(name string, value float64, tags ...string) (err error) {
	for _, collector := range collectors {
		err = collector.Histogram(name, value, tags...)
		if err != nil {
			return
		}
	}
	return
}

// Distribution sets a distribution value and writes to the collectors.
func (collectors MultiCollector) Distribution(name string, value float64, tags ...string) (err error) {
	for _, collector := range collectors {
		err = collector.Distribution(name, value, tags...)
		if err != nil {
			return
		}
	}
	return
}

// TimeInMilliseconds sets a timing value and writes to the different hosts
func (collectors MultiCollector) TimeInMilliseconds(name string, value time.Duration, tags ...string) (err error) {
	for _, collector := range collectors {
		err = collector.TimeInMilliseconds(name, value, tags...)
		if err != nil {
			return
		}
	}
	return
}

// Flush forces a flush on all collectors.
func (collectors MultiCollector) Flush() (err error) {
	for _, collector := range collectors {
		err = collector.Flush()
		if err != nil {
			return
		}
	}
	return
}

// Close closes all collectors.
func (collectors MultiCollector) Close() (err error) {
	for _, collector := range collectors {
		err = collector.Close()
		if err != nil {
			return
		}
	}
	return
}
