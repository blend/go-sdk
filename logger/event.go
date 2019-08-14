package logger

import (
	"io"
	"time"
)

// Event is an interface representing methods necessary to trigger listeners.
type Event interface {
	GetFlag() string
}

// TimestampProvider is a type that provides a timestamp.
type TimestampProvider interface {
	Timestamp() time.Time
}

// TextWritable is an event that can be written.
type TextWritable interface {
	WriteText(TextFormatter, io.Writer)
}

// JSONWritable is a type that implements a decompose method.
// This is used by the json serializer.
type JSONWritable interface {
	Decompose() map[string]interface{}
}
