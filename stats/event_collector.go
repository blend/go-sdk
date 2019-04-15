package stats

import (
	"time"

	"github.com/blend/go-sdk/ex"
)

// EventCollector is a collector for events.
type EventCollector interface {
	Taggable
	SendEvent(Event) error
	CreateEvent(title, text string, tags ...string) Event
}

const (
	// EventAlertTypeInfo is the "info" AlertType for events
	EventAlertTypeInfo = "info"
	// EventAlertTypeError is the "error" AlertType for events
	EventAlertTypeError = "error"
	// EventAlertTypeWarning is the "warning" AlertType for events
	EventAlertTypeWarning = "warning"
	// EventAlertTypeSuccess is the "success" AlertType for events
	EventAlertTypeSuccess = "success"
)

const (
	// EventPriorityNormal is the "normal" Priority for events.
	EventPriorityNormal = "normal"
	// EventPriorityLow is the "low" Priority for events.
	EventPriorityLow = "low"
)

// Event is an event to be collected.
type Event struct {
	// Title of the event.  Required.
	Title string
	// Text is the description of the event.  Required.
	Text string
	// Timestamp is a timestamp for the event.  If not provided, the dogstatsd
	// server will set this to the current time.
	Timestamp time.Time
	// Hostname for the event.
	Hostname string
	// AggregationKey groups this event with others of the same key.
	AggregationKey string
	// Priority of the event.  Can be statsd.Low or statsd.Normal.
	Priority string
	// SourceTypeName is a source type for the event.
	SourceTypeName string
	// AlertType can be statsd.Info, statsd.Error, statsd.Warning, or statsd.Success.
	// If absent, the default value applied by the dogstatsd server is Info.
	AlertType string
	// Tags for the event.
	Tags []string
}

// Check verifies that an event is valid.
func (e Event) Check() error {
	if len(e.Title) == 0 {
		return ex.Class("event title is required")
	}
	if len(e.Text) == 0 {
		return ex.Class("event text is required")
	}
	return nil
}
