package migration

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/logger"
)

var (
	_ logger.Event        = (*Event)(nil)
	_ logger.TextWritable = (*Event)(nil)
	_ json.Marshaler      = (*Event)(nil)
)

const (
	// Flag is a logger event flag.
	Flag = "db.migration"
	// FlagStats is a logger event flag.
	FlagStats = "db.migration.stats"
)

// NewEvent returns a new event.
func NewEvent(result, body string, labels ...string) *Event {
	return &Event{
		EventMeta: logger.NewEventMeta(Flag),
		Result:    result,
		Body:      body,
		Labels:    labels,
	}
}

// Event is a migration logger event.
type Event struct {
	*logger.EventMeta

	Result string
	Body   string
	Labels []string
}

// WriteText writes the migration event as text.
func (e Event) WriteText(tf logger.TextFormatter, wr io.Writer) {
	resultColor := ansi.ColorBlue
	switch e.Result {
	case "skipped":
		resultColor = ansi.ColorYellow
	case "failed":
		resultColor = ansi.ColorRed
	}

	if len(e.Result) > 0 {
		io.WriteString(wr, tf.Colorize("--", ansi.ColorLightBlack))
		io.WriteString(wr, logger.Space)
		io.WriteString(wr, tf.Colorize(e.Result, resultColor))
	}

	if len(e.Labels) > 0 {
		io.WriteString(wr, logger.Space)
		io.WriteString(wr, strings.Join(e.Labels, " > "))
	}

	if len(e.Body) > 0 {
		io.WriteString(wr, logger.Space)
		io.WriteString(wr, tf.Colorize("--", ansi.ColorLightBlack))
		io.WriteString(wr, logger.Space)
		io.WriteString(wr, e.Body)
	}
}

// MarshalJSON implements json.Marshaler.
func (e Event) MarshalJSON() ([]byte, error) {
	return json.Marshal(logger.MergeDecomposed(e.EventMeta.Decompose(), map[string]interface{}{
		"result": e.Result,
		"labels": e.Labels,
		"body":   e.Body,
	}))
}
