/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package grpcutil

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/timeutil"
)

// Logger flags
const (
	FlagRPCStreamMessage = "rpc.stream.message"
)

// these are compile time assertions
var (
	_	logger.Event		= (*RPCStreamMessageEvent)(nil)
	_	logger.TextWritable	= (*RPCStreamMessageEvent)(nil)
	_	logger.JSONWritable	= (*RPCStreamMessageEvent)(nil)
)

// StreamMessageDirection is the direction the message was sent.
type StreamMessageDirection string

// constants
const (
	StreamMessageDirectionReceive	= "recv"
	StreamMessageDirectionSend	= "send"
)

// NewRPCStreamMessageEvent creates a new rpc stream message event.
func NewRPCStreamMessageEvent(method string, direction StreamMessageDirection, elapsed time.Duration, options ...RPCStreamMessageEventOption) RPCStreamMessageEvent {
	rpe := RPCStreamMessageEvent{
		RPCEvent: RPCEvent{
			Engine:		EngineGRPC,
			Method:		method,
			Elapsed:	elapsed,
		},
		Direction:	direction,
	}
	for _, opt := range options {
		opt(&rpe)
	}
	return rpe
}

// NewRPCStreamMessageEventListener returns a new rpc stream message event event listener.
func NewRPCStreamMessageEventListener(listener func(context.Context, RPCStreamMessageEvent)) logger.Listener {
	return func(ctx context.Context, e logger.Event) {
		if typed, isTyped := e.(RPCStreamMessageEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// NewRPCStreamMessageEventFilter returns a new rpc stream message event filter.
func NewRPCStreamMessageEventFilter(filter func(context.Context, RPCStreamMessageEvent) (RPCStreamMessageEvent, bool)) logger.Filter {
	return func(ctx context.Context, e logger.Event) (logger.Event, bool) {
		if typed, isTyped := e.(RPCStreamMessageEvent); isTyped {
			return filter(ctx, typed)
		}
		return e, false
	}
}

// RPCStreamMessageEventOption is a mutator for RPCEvents.
type RPCStreamMessageEventOption func(*RPCStreamMessageEvent)

// OptRPCStreamMessageEngine sets a field on the event.
func OptRPCStreamMessageEngine(value string) RPCStreamMessageEventOption {
	return func(e *RPCStreamMessageEvent) { e.Engine = value }
}

// OptRPCStreamMessagePeer sets a field on the event.
func OptRPCStreamMessagePeer(value string) RPCStreamMessageEventOption {
	return func(e *RPCStreamMessageEvent) { e.Peer = value }
}

// OptRPCStreamMessageMethod sets a field on the event.
func OptRPCStreamMessageMethod(value string) RPCStreamMessageEventOption {
	return func(e *RPCStreamMessageEvent) { e.Method = value }
}

// OptRPCStreamMessageDirection sets a field on the event.
func OptRPCStreamMessageDirection(value StreamMessageDirection) RPCStreamMessageEventOption {
	return func(e *RPCStreamMessageEvent) { e.Direction = value }
}

// OptRPCStreamMessageUserAgent sets a field on the event.
func OptRPCStreamMessageUserAgent(value string) RPCStreamMessageEventOption {
	return func(e *RPCStreamMessageEvent) { e.UserAgent = value }
}

// OptRPCStreamMessageAuthority sets a field on the event.
func OptRPCStreamMessageAuthority(value string) RPCStreamMessageEventOption {
	return func(e *RPCStreamMessageEvent) { e.Authority = value }
}

// OptRPCStreamMessageContentType sets a field on the event.
func OptRPCStreamMessageContentType(value string) RPCStreamMessageEventOption {
	return func(e *RPCStreamMessageEvent) { e.ContentType = value }
}

// OptRPCStreamMessageElapsed sets a field on the event.
func OptRPCStreamMessageElapsed(value time.Duration) RPCStreamMessageEventOption {
	return func(e *RPCStreamMessageEvent) { e.Elapsed = value }
}

// OptRPCStreamMessageErr sets a field on the event.
func OptRPCStreamMessageErr(value error) RPCStreamMessageEventOption {
	return func(e *RPCStreamMessageEvent) { e.Err = value }
}

// RPCStreamMessageEvent is an event type for rpc
type RPCStreamMessageEvent struct {
	RPCEvent
	Direction	StreamMessageDirection
}

// GetFlag implements Event.
func (e RPCStreamMessageEvent) GetFlag() string	{ return FlagRPCStreamMessage }

// WriteText implements TextWritable.
func (e RPCStreamMessageEvent) WriteText(tf logger.TextFormatter, wr io.Writer) {
	if e.Engine != "" {
		fmt.Fprint(wr, "[")
		fmt.Fprint(wr, tf.Colorize(e.Engine, ansi.ColorLightWhite))
		fmt.Fprint(wr, "]")
	}
	if e.Method != "" {
		if e.Engine != "" {
			fmt.Fprint(wr, logger.Space)
		}
		fmt.Fprint(wr, tf.Colorize(e.Method, ansi.ColorBlue))
	}
	if e.Direction != "" {
		fmt.Fprint(wr, logger.Space)
		fmt.Fprint(wr, e.Direction)
	}
	if e.Peer != "" {
		fmt.Fprint(wr, logger.Space)
		fmt.Fprint(wr, e.Peer)
	}
	if e.Authority != "" {
		fmt.Fprint(wr, logger.Space)
		fmt.Fprint(wr, e.Authority)
	}
	if e.UserAgent != "" {
		fmt.Fprint(wr, logger.Space)
		fmt.Fprint(wr, e.UserAgent)
	}
	if e.ContentType != "" {
		fmt.Fprint(wr, logger.Space)
		fmt.Fprint(wr, e.ContentType)
	}

	fmt.Fprint(wr, logger.Space)
	fmt.Fprint(wr, e.Elapsed.String())

	if e.Err != nil {
		fmt.Fprint(wr, logger.Space)
		fmt.Fprint(wr, tf.Colorize("failed", ansi.ColorRed))
	}
}

// Decompose implements JSONWritable.
func (e RPCStreamMessageEvent) Decompose() map[string]interface{} {
	return map[string]interface{}{
		"engine":	e.Engine,
		"peer":		e.Peer,
		"method":	e.Method,
		"direction":	e.Direction,
		"userAgent":	e.UserAgent,
		"authority":	e.Authority,
		"contentType":	e.ContentType,
		"elapsed":	timeutil.Milliseconds(e.Elapsed),
		"err":		e.Err,
	}
}
