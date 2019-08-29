package cron

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/timeutil"
)

// these are compile time assertions
var (
	_ logger.Event        = (*Event)(nil)
	_ logger.TextWritable = (*Event)(nil)
	_ logger.JSONWritable = (*Event)(nil)
)

// NewEventListener returns a new event listener.
func NewEventListener(listener func(context.Context, Event)) logger.Listener {
	return func(ctx context.Context, e logger.Event) {
		if typed, isTyped := e.(Event); isTyped {
			listener(ctx, typed)
		}
	}
}

// NewEvent creates a new event with a given set of optional options.
func NewEvent(flag, jobName string, options ...EventOption) Event {
	e := Event{
		Flag:    flag,
		JobName: jobName,
	}

	for _, option := range options {
		option(&e)
	}
	return e
}

// EventOption is an option for an Event.
type EventOption func(*Event)

// OptEventEnabled sets an enabled provider.
func OptEventEnabled(enabled bool) EventOption {
	return func(e *Event) {
		e.EnabledProvider = func() bool { return enabled }
	}
}

// OptEventWritable sets a writable provider.
func OptEventWritable(enabled bool) EventOption {
	return func(e *Event) {
		e.EnabledProvider = func() bool { return enabled }
	}
}

// OptEventJobInvocation sets a field.
func OptEventJobInvocation(jobInvocation string) EventOption {
	return func(e *Event) { e.JobInvocation = jobInvocation }
}

// OptEventErr sets a field.
func OptEventErr(err error) EventOption {
	return func(e *Event) { e.Err = err }
}

// OptEventElapsed sets a field.
func OptEventElapsed(elapsed time.Duration) EventOption {
	return func(e *Event) { e.Elapsed = elapsed }
}

// Event is an event.
type Event struct {
	Flag string

	EnabledProvider  func() bool
	WritableProvider func() bool

	JobName       string
	JobInvocation string
	Err           error
	Elapsed       time.Duration
}

// GetFlag implements logger.Event.
func (e Event) GetFlag() string { return e.Flag }

// Complete returns if the event completed.
func (e Event) Complete() bool {
	return e.Flag == FlagComplete
}

// IsEnabled is a
func (e Event) IsEnabled() bool {
	if e.EnabledProvider != nil {
		return e.EnabledProvider()
	}
	return true
}

// IsWritable is a logger interface to disable writing the events.
func (e Event) IsWritable() bool {
	if e.WritableProvider != nil {
		return e.WritableProvider()
	}
	return true
}

// WriteText implements logger.TextWritable.
func (e Event) WriteText(tf logger.TextFormatter, wr io.Writer) {
	if e.JobInvocation != "" {
		io.WriteString(wr, fmt.Sprintf("[%s > %s]", tf.Colorize(e.JobName, ansi.ColorBlue), tf.Colorize(e.JobInvocation, ansi.ColorBlue)))
	} else {
		io.WriteString(wr, fmt.Sprintf("[%s]", tf.Colorize(e.JobName, ansi.ColorBlue)))
	}

	if e.Elapsed > 0 {
		io.WriteString(wr, logger.Space)
		io.WriteString(wr, fmt.Sprintf("(%v)", e.Elapsed))
	}
}

// Decompose implements logger.JSONWritable.
func (e Event) Decompose() map[string]interface{} {
	return map[string]interface{}{
		"jobName": e.JobName,
		"err":     e.Err,
		"elapsed": timeutil.Milliseconds(e.Elapsed),
	}
}
