package logger

import (
	"time"
)

// these are compile time assertions
var (
	_ Event            = &MessageEvent{}
	_ EventHeadings    = &MessageEvent{}
	_ EventLabels      = &MessageEvent{}
	_ EventAnnotations = &MessageEvent{}
)

// NewFrontendActionEvent creates a new query event.
func NewFrontendActionEvent(body []byte) *FrontendActionEvent {
	return &FrontendActionEvent{
		EventMeta: NewEventMeta(FrontendAction),
		body:      body,
	}
}

// FrontendActionEvent represents an event on the frontend
type FrontendActionEvent struct {
	*EventMeta

	body []byte
}

// Body returns the event body.
func (f *FrontendActionEvent) Body() []byte {
	return f.body
}

// Flag returns the event flag
func (f *FrontendActionEvent) Flag() Flag {
	return f.flag
}

// Timestamp returns the event timestamp
func (f *FrontendActionEvent) Timestamp() time.Time {
	return f.ts
}

// String returns the event as a string
func (f *FrontendActionEvent) String() string {
	return string(f.body)
}

// WithBody sets the body.
func (f *FrontendActionEvent) WithBody(body []byte) *FrontendActionEvent {
	f.body = body
	return f
}

// WithTimestamp sets the timestamp.
func (f *FrontendActionEvent) WithTimestamp(ts time.Time) *FrontendActionEvent {
	f.ts = ts
	return f
}

// WithLabel sets a label on the event for later filtering.
func (f *FrontendActionEvent) WithLabel(key, value string) *FrontendActionEvent {
	f.AddLabelValue(key, value)
	return f
}

// WithAnnotation adds an annotation to the event.
func (f *FrontendActionEvent) WithAnnotation(key, value string) *FrontendActionEvent {
	f.AddAnnotationValue(key, value)
	return f
}

// WithFlag sets the flag.
func (f *FrontendActionEvent) WithFlag(flag Flag) *FrontendActionEvent {
	f.flag = flag
	return f
}

// Labels returns the event's labels
func (f *FrontendActionEvent) Labels() map[string]string {
	return f.labels
}

// Annotations returns the event's annotations
func (f *FrontendActionEvent) Annotations() map[string]string {
	return f.annotations
}
