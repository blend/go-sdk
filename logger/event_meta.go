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

// EventMeta is the metadata common to events.
// It is useful for ensuring you have the minimum required fields on your events, and its typically embedded in types.
type EventMeta struct {
	Flag      string
	Timestamp time.Time
	FlagColor ansi.Color
	Fields    map[string]string
}

// GetFlag returns the event flag.
func (em EventMeta) GetFlag() string { return em.Flag }

// GetTimestamp returns the event timestamp.
func (em EventMeta) GetTimestamp() time.Time { return em.Timestamp }

// GetFlagColor returns the event flag color
func (em EventMeta) GetFlagColor() ansi.Color { return em.FlagColor }

// Decompose decomposes the object into a map[string]interface{}.
func (em EventMeta) Decompose() map[string]interface{} {
	output := map[string]interface{}{
		FieldFlag:      em.Flag,
		FieldTimestamp: em.Timestamp.Format(time.RFC3339Nano),
		FieldFields:    em.Fields,
	}
	return output
}
