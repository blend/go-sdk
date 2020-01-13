package secrets

import (
	"io"
	"net/http"
	"strings"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/logger"
)

var (
	_ logger.Event        = (*Event)(nil)
	_ logger.TextWritable = (*Event)(nil)
	_ logger.JSONWritable = (*Event)(nil)
)

const (
	// Flag is the logger flag.
	Flag = "secrets"
)

// NewEvent returns a new event from a request.
func NewEvent(req *http.Request) *Event {
	return &Event{
		Remote: req.URL.Host,
		Method: req.Method,
		Key:    strings.TrimPrefix(req.URL.Path, "/v1/"),
	}
}

// Event is an event.
type Event struct {
	Remote string
	Method string
	Key    string
}

// GetFlag implements logger.Event.
func (e Event) GetFlag() string { return Flag }

// WriteText writes text for the event.
func (e *Event) WriteText(tf logger.TextFormatter, wr io.Writer) {
	io.WriteString(wr, "["+tf.Colorize(e.Method, ansi.ColorBlue)+"]")
	io.WriteString(wr, logger.Space)
	io.WriteString(wr, e.Remote)
	io.WriteString(wr, logger.Space)
	io.WriteString(wr, e.Key)
}

// Decompose impements logger.JSONWritable.
func (e *Event) Decompose() map[string]interface{} {
	return map[string]interface{}{
		"remote": e.Remote,
		"method": e.Method,
		"key":    e.Key,
	}
}
