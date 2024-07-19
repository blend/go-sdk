/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package grpcutil

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/timeutil"
)

// Logger flags
const (
	FlagRPC = "rpc"
)

// these are compile time assertions
var (
	_ logger.Event        = (*RPCEvent)(nil)
	_ logger.TextWritable = (*RPCEvent)(nil)
	_ logger.JSONWritable = (*RPCEvent)(nil)
)

// NewRPCEvent creates a new rpc event.
func NewRPCEvent(method string, elapsed time.Duration, options ...RPCEventOption) RPCEvent {
	rpe := RPCEvent{
		Engine:  EngineGRPC,
		Method:  method,
		Elapsed: elapsed,
	}
	for _, opt := range options {
		opt(&rpe)
	}
	return rpe
}

// NewRPCEventListener returns a new web request event listener.
func NewRPCEventListener(listener func(context.Context, RPCEvent)) logger.Listener {
	return func(ctx context.Context, e logger.Event) {
		if typed, isTyped := e.(RPCEvent); isTyped {
			listener(ctx, typed)
		}
	}
}

// NewRPCEventFilter returns a new rpc event filter.
func NewRPCEventFilter(filter func(context.Context, RPCEvent) (RPCEvent, bool)) logger.Filter {
	return func(ctx context.Context, e logger.Event) (logger.Event, bool) {
		if typed, isTyped := e.(RPCEvent); isTyped {
			return filter(ctx, typed)
		}
		return e, false
	}
}

// RPCEventOption is a mutator for RPCEvents.
type RPCEventOption func(*RPCEvent)

// OptRPCEngine sets a field on the event.
func OptRPCEngine(value string) RPCEventOption {
	return func(e *RPCEvent) { e.Engine = value }
}

// OptRPCPeer sets a field on the event.
func OptRPCPeer(value string) RPCEventOption {
	return func(e *RPCEvent) { e.Peer = value }
}

// OptRPCMethod sets a field on the event.
func OptRPCMethod(value string) RPCEventOption {
	return func(e *RPCEvent) { e.Method = value }
}

// OptRPCUserAgent sets a field on the event.
func OptRPCUserAgent(value string) RPCEventOption {
	return func(e *RPCEvent) { e.UserAgent = value }
}

// OptRPCAuthority sets a field on the event.
func OptRPCAuthority(value string) RPCEventOption {
	return func(e *RPCEvent) { e.Authority = value }
}

// OptRPCContentType sets a field on the event.
func OptRPCContentType(value string) RPCEventOption {
	return func(e *RPCEvent) { e.ContentType = value }
}

// OptRPCElapsed sets a field on the event.
func OptRPCElapsed(value time.Duration) RPCEventOption {
	return func(e *RPCEvent) { e.Elapsed = value }
}

// OptRPCErr sets a field on the event.
func OptRPCErr(value error) RPCEventOption {
	return func(e *RPCEvent) { e.Err = value }
}

// RPCEvent is an event type for rpc
type RPCEvent struct {
	Engine      string
	Peer        string
	Method      string
	UserAgent   string
	Authority   string
	ContentType string
	Elapsed     time.Duration
	Err         error
}

// GetFlag implements Event.
func (e RPCEvent) GetFlag() string { return FlagRPC }

// WriteText implements TextWritable.
func (e RPCEvent) WriteText(tf logger.TextFormatter, wr io.Writer) {
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

		if s, ok := status.FromError(e.Err); ok {
			fmt.Fprint(wr, tf.Colorize(fmt.Sprintf("failed (%[1]d - %[1]v)", s.Code()), ansi.ColorRed))
		} else {
			fmt.Fprint(wr, tf.Colorize("failed", ansi.ColorRed))
		}
	}
}

// Decompose implements JSONWritable.
func (e RPCEvent) Decompose() map[string]interface{} {
	var code codes.Code
	if s, ok := status.FromError(e.Err); ok {
		code = s.Code()
	}
	return map[string]interface{}{
		"engine":      e.Engine,
		"peer":        e.Peer,
		"method":      e.Method,
		"userAgent":   e.UserAgent,
		"authority":   e.Authority,
		"contentType": e.ContentType,
		"elapsed":     timeutil.Milliseconds(e.Elapsed),
		"err":         e.Err,
		"code":        code,
	}
}
