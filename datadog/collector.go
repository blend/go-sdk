package datadog

import (
	"fmt"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/statsd"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/timeutil"
)

// Assert that the datadog collector implements stats.Collector and stats.EventCollector.
var (
	_ stats.Collector      = (*Collector)(nil)
	_ stats.EventCollector = (*Collector)(nil)
)

// New returns a new stats collector from a config.
func New(cfg Config, opts ...statsd.Option) (*Collector, error) {
	var client *statsd.Client
	var err error
	if cfg.BufferedOrDefault() {
		opts = append(opts, statsd.WithMaxMessagesPerPayload(cfg.BufferDepthOrDefault()))
	}
	client, err = statsd.New(cfg.GetAddress(), opts...)
	if err != nil {
		return nil, err
	}
	if len(cfg.Namespace) > 0 {
		client.Namespace = strings.ToLower(cfg.Namespace) + "."
	}
	collector := &Collector{
		client:      client,
		defaultTags: cfg.DefaultTags,
	}

	collector.AddDefaultTag("service", env.Env().String(env.VarServiceName))
	collector.AddDefaultTag("env", env.Env().String(env.VarServiceEnv))
	collector.AddDefaultTag("hostname", env.Env().String(env.VarHostname))

	return collector, nil
}

// MustNew returns a new stats collector from a config, but panics on error.
func MustNew(cfg Config) *Collector {
	collector, err := New(cfg)
	if err != nil {
		panic(err)
	}
	return collector
}

// Collector is a class that wraps the statsd collector we're using.
type Collector struct {
	client      *statsd.Client
	defaultTags []string
}

// AddDefaultTag adds a new default tag and returns a reference to the collector.
func (dc *Collector) AddDefaultTag(key, value string) {
	dc.defaultTags = append(dc.defaultTags, fmt.Sprintf("%s:%s", key, value))
}

// DefaultTags returns the default tags for the collector.
func (dc *Collector) DefaultTags() []string {
	return dc.defaultTags
}

// Count increments a counter by a value.
func (dc *Collector) Count(name string, value int64, tags ...string) error {
	return dc.client.Count(name, value, dc.tags(tags...), 1.0)
}

// Increment increments a counter by 1.
func (dc *Collector) Increment(name string, tags ...string) error {
	return dc.client.Count(name, 1, dc.tags(tags...), 1.0)
}

// Gauge sets a gauge value.
func (dc *Collector) Gauge(name string, value float64, tags ...string) error {
	return dc.client.Gauge(name, value, dc.tags(tags...), 1.0)
}

// Histogram sets a guage value.
func (dc *Collector) Histogram(name string, value float64, tags ...string) error {
	return dc.client.Histogram(name, value, dc.tags(tags...), 1.0)
}

// TimeInMilliseconds sets a timing value.
func (dc *Collector) TimeInMilliseconds(name string, value time.Duration, tags ...string) error {
	return dc.client.TimeInMilliseconds(name, timeutil.Milliseconds(value), dc.tags(tags...), 1.0)
}

// Flush forces a flush of all the queued statsd payloads.
func (dc *Collector) Flush() error {
	if dc.client == nil {
		return nil
	}
	return dc.client.Flush()
}

// SimpleEvent sends an event w/ title and text
func (dc *Collector) SimpleEvent(title, text string) error {
	return dc.client.SimpleEvent(title, text)
}

// SendEvent sends any *statsd.Event
func (dc *Collector) SendEvent(event stats.Event) error {
	return dc.client.Event(ConvertEvent(event))
}

// CreateEvent makes a new Event with the collectors default tags.
func (dc *Collector) CreateEvent(title, text string, tags ...string) stats.Event {
	return stats.Event{
		Title: title,
		Text:  text,
		Tags:  dc.tags(tags...),
	}
}

// helpers
func (dc *Collector) tags(tags ...string) []string {
	return append(dc.defaultTags, tags...)
}

// ConvertEvent converts a stats event to a statsd (datadog) event.
func ConvertEvent(e stats.Event) *statsd.Event {
	return &statsd.Event{
		Title:          e.Title,
		Text:           e.Text,
		Timestamp:      e.Timestamp,
		Hostname:       e.Hostname,
		AggregationKey: e.AggregationKey,
		Priority:       statsd.EventPriority(e.Priority),
		SourceTypeName: e.SourceTypeName,
		AlertType:      statsd.EventAlertType(e.AlertType),
		Tags:           e.Tags,
	}
}
