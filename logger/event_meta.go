package logger

import (
	"time"

	"github.com/blend/go-sdk/ansi"
)

// these are compile time assertions
var (
	_ Event = (*EventMeta)(nil)
)

// NewEventMeta returns a new event meta.
func NewEventMeta(flag string, options ...EventMetaOption) *EventMeta {
	em := &EventMeta{
		Flag:      flag,
		Timestamp: time.Now().UTC(),
	}
	for _, option := range options {
		option(em)
	}
	return em
}

// EventMetaOption is an option for event metas.
type EventMetaOption func(*EventMeta)

// OptEventMetaFlag sets the event flag.
func OptEventMetaFlag(flag string) EventMetaOption {
	return func(em *EventMeta) { em.Flag = flag }
}

// OptEventMetaTimestamp sets the event timestamp.
func OptEventMetaTimestamp(ts time.Time) EventMetaOption {
	return func(em *EventMeta) { em.Timestamp = ts }
}

// OptEventMetaFlagColor sets the event flag color.
func OptEventMetaFlagColor(color ansi.Color) EventMetaOption {
	return func(em *EventMeta) { em.FlagColor = color }
}

// OptEventMetaLabels sets the event meta labels.
func OptEventMetaLabels(labels Labels) EventMetaOption {
	return func(em *EventMeta) { em.Labels = labels }
}

// OptEventMetaAnnotations sets the event meta labels.
func OptEventMetaAnnotations(annotations Annotations) EventMetaOption {
	return func(em *EventMeta) { em.Annotations = annotations }
}

// EventMeta is the metadata common to events.
// It is useful for ensuring you have the minimum required fields on your events, and its typically embedded in types.
type EventMeta struct {
	Labels
	Annotations

	Flag      string
	Timestamp time.Time
	FlagColor ansi.Color
}

// GetLabels returns the labels.
func (em EventMeta) GetLabels() Labels { return em.Labels }

// GetAnnotations returns the annotations.
func (em EventMeta) GetAnnotations() Annotations { return em.Annotations }

// GetFlag returns the event flag.
func (em EventMeta) GetFlag() string { return em.Flag }

// GetTimestamp returns the event timestamp.
func (em EventMeta) GetTimestamp() time.Time { return em.Timestamp }

// GetFlagColor returns the event flag color
func (em EventMeta) GetFlagColor() ansi.Color { return em.FlagColor }

// Decompose decomposes the object into a map[string]interface{}.
func (em EventMeta) Decompose() map[string]interface{} {
	output := map[string]interface{}{
		FieldFlag:        em.Flag,
		FieldTimestamp:   em.Timestamp.Format(time.RFC3339Nano),
		FieldLabels:      em.Labels.Decompose(),
		FieldAnnotations: em.Annotations.Decompose(),
	}
	return output
}
