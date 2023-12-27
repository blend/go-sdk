/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package migration

import (
	"fmt"
	"io"
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
	// Flag is a logger event flag.
	Flag = "db.migration"
	// FlagStats is a logger event flag.
	FlagStats = "db.migration.stats"
)

// NewEvent returns a new event.
func NewEvent(result, body string, labels ...string) *Event {
	return &Event{
		Result: result,
		Body:   body,
		Labels: labels,
	}
}

// Event is a migration logger event.
type Event struct {
	Result string
	Body   string
	Labels []string
}

// GetFlag implements logger.Event.
func (e Event) GetFlag() string { return Flag }

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
		fmt.Fprint(wr, tf.Colorize("--", ansi.ColorLightBlack))
		fmt.Fprint(wr, logger.Space)
		fmt.Fprint(wr, tf.Colorize(e.Result, resultColor))
	}

	if len(e.Labels) > 0 {
		fmt.Fprint(wr, logger.Space)
		fmt.Fprint(wr, strings.Join(e.Labels, " > "))
	}

	if len(e.Body) > 0 {
		fmt.Fprint(wr, logger.Space)
		fmt.Fprint(wr, tf.Colorize("--", ansi.ColorLightBlack))
		fmt.Fprint(wr, logger.Space)
		fmt.Fprint(wr, e.Body)
	}
}

// Decompose implements logger.JSONWritable.
func (e Event) Decompose() map[string]interface{} {
	return map[string]interface{}{
		"result": e.Result,
		"labels": e.Labels,
		"body":   e.Body,
	}
}

var (
	_ logger.Event        = (*StatsEvent)(nil)
	_ logger.TextWritable = (*StatsEvent)(nil)
	_ logger.JSONWritable = (*StatsEvent)(nil)
)

// NewStatsEvent returns a new stats event.
func NewStatsEvent(applied, skipped, failed, total int) *StatsEvent {
	return &StatsEvent{
		applied: applied,
		skipped: skipped,
		failed:  failed,
		total:   total,
	}
}

// StatsEvent is a migration logger event.
type StatsEvent struct {
	applied int
	skipped int
	failed  int
	total   int
}

// GetFlag implements logger.Event.
func (se StatsEvent) GetFlag() string { return FlagStats }

// WriteText writes the event to a text writer.
func (se StatsEvent) WriteText(tf logger.TextFormatter, wr io.Writer) {
	fmt.Fprintf(wr, "%s applied %s skipped %s failed %s total",
		tf.Colorize(fmt.Sprintf("%d", se.applied), ansi.ColorGreen),
		tf.Colorize(fmt.Sprintf("%d", se.skipped), ansi.ColorLightGreen),
		tf.Colorize(fmt.Sprintf("%d", se.failed), ansi.ColorRed),
		tf.Colorize(fmt.Sprintf("%d", se.total), ansi.ColorLightWhite),
	)
}

// Decompose implements logger.JSONWritable.
func (se StatsEvent) Decompose() map[string]interface{} {
	return map[string]interface{}{
		StatApplied: se.applied,
		StatSkipped: se.skipped,
		StatFailed:  se.failed,
		StatTotal:   se.total,
	}
}
