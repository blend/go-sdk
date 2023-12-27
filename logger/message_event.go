/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package logger

import (
	"context"
	"fmt"
	"io"
	"time"
)

// these are compile time assertions
var (
	_ Event = (*MessageEvent)(nil)
)

// NewMessageEvent returns a new message event.
func NewMessageEvent(flag, text string, options ...MessageEventOption) MessageEvent {
	me := MessageEvent{
		Flag:	flag,
		Text:	text,
	}
	for _, opt := range options {
		opt(&me)
	}
	return me
}

// NewMessageEventListener returns a new message event listener.
func NewMessageEventListener(listener func(context.Context, MessageEvent)) Listener {
	return func(ctx context.Context, e Event) {
		if typed, isTyped := e.(MessageEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// NewMessageEventFilter returns a new message event filter.
func NewMessageEventFilter(filter func(context.Context, MessageEvent) (MessageEvent, bool)) Filter {
	return func(ctx context.Context, e Event) (Event, bool) {
		if typed, isTyped := e.(MessageEvent); isTyped {
			return filter(ctx, typed)
		}
		return e, false
	}
}

// MessageEventOption mutates a message event.
type MessageEventOption func(*MessageEvent)

// OptMessageFlag sets a field on a message event.
func OptMessageFlag(flag string) MessageEventOption {
	return func(me *MessageEvent) { me.Flag = flag }
}

// OptMessageText sets a field on a message event.
func OptMessageText(text string) MessageEventOption {
	return func(me *MessageEvent) { me.Text = text }
}

// OptMessageElapsed sets a field on a message event.
func OptMessageElapsed(elapsed time.Duration) MessageEventOption {
	return func(me *MessageEvent) { me.Elapsed = elapsed }
}

// OptMessageRestricted sets a field on a message event.
func OptMessageRestricted(restricted bool) MessageEventOption {
	return func(me *MessageEvent) { me.Restricted = restricted }
}

// MessageEvent is a common type of message.
type MessageEvent struct {
	Flag		string
	Text		string
	Restricted	bool
	Elapsed		time.Duration
}

// GetFlag implements Event.
func (e MessageEvent) GetFlag() string	{ return e.Flag }

// WriteText implements TextWritable.
func (e MessageEvent) WriteText(formatter TextFormatter, output io.Writer) {
	fmt.Fprint(output, e.Text)
	if e.Elapsed > 0 {
		fmt.Fprint(output, Space)
		fmt.Fprint(output, "("+e.Elapsed.String()+")")
	}
}

// Decompose implements json.Marshaler.
func (e MessageEvent) Decompose() map[string]interface{} {
	m := map[string]interface{}{
		FieldText:		e.Text,
		FieldRestricted:	e.Restricted,
	}
	if e.Elapsed > 0 {
		m[FieldElapsed] = e.Elapsed
	}
	return m
}
