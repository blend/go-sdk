/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package grpcutil

import (
	"context"

	"google.golang.org/grpc"
)

// Tracer is the full tracer.
type Tracer interface {
	ServerTracer
	ClientTracer
}

// ServerTracer is a type that starts traces.
type ServerTracer interface {
	StartServerUnary(ctx context.Context, method string) (context.Context, TraceFinisher, error)
	StartServerStream(ctx context.Context, method string) (context.Context, TraceFinisher, error)
}

// ClientTracer is a type that starts traces.
type ClientTracer interface {
	StartClientUnary(ctx context.Context, remoteAddr, method string) (context.Context, TraceFinisher, error)
	StartClientStream(ctx context.Context, remoteAddr, method string) (context.Context, TraceFinisher, error)
}

// TraceFinisher is a finisher for traces
type TraceFinisher interface {
	Finish(err error)
}

// TracedServerUnary returns a unary server interceptor.
func TracedServerUnary(tracer ServerTracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, args interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (result interface{}, err error) {
		if tracer == nil {
			return handler(ctx, args)
		}
		var finisher TraceFinisher
		ctx, finisher, err = tracer.StartServerUnary(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		defer func() {
			finisher.Finish(err)
		}()
		result, err = handler(ctx, args)
		return
	}
}

// TracedServerStream returns a grpc streaming interceptor.
func TracedServerStream(tracer ServerTracer) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if tracer == nil {
			return handler(srv, ss)
		}
		var finisher TraceFinisher
		var err error
		var ctx context.Context
		ctx, finisher, err = tracer.StartServerStream(ss.Context(), info.FullMethod)
		if err != nil {
			return err
		}
		defer func() {
			finisher.Finish(err)
		}()
		err = handler(srv, &contextServerStream{ServerStream: ss, ctx: ctx})
		return err
	}
}

// spanServerStream wraps around the embedded grpc.ServerStream, and
// intercepts calls to `Context()` returning a context with the span information injected.
//
// NOTE: you can extend this type to intercept calls to `SendMsg` and `RecvMsg` if you want to
// add tracing handling for individual stream calls.
type contextServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (cs *contextServerStream) Context() context.Context {
	if cs.ctx != nil {
		return cs.ctx
	}
	return cs.ServerStream.Context()
}

// TracedClientUnary implements the unary client interceptor based on a tracer.
func TracedClientUnary(tracer ClientTracer) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		if tracer == nil {
			err = invoker(ctx, method, req, reply, cc, opts...)
			return
		}
		var finisher TraceFinisher
		ctx, finisher, err = tracer.StartClientUnary(ctx, cc.Target(), method)
		if err != nil {
			return
		}
		defer func() {
			finisher.Finish(err)
		}()
		err = invoker(ctx, method, req, reply, cc, opts...)
		return
	}
}

// TracedClientStream implements the stream client interceptor based on a tracer.
func TracedClientStream(tracer ClientTracer) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (cs grpc.ClientStream, err error) {
		if tracer == nil {
			cs, err = streamer(ctx, desc, cc, method, opts...)
			return
		}
		var finisher TraceFinisher
		ctx, finisher, err = tracer.StartClientStream(ctx, cc.Target(), method)
		if err != nil {
			return
		}
		defer func() {
			finisher.Finish(err)
		}()
		cs, err = streamer(ctx, desc, cc, method, opts...)
		return
	}
}
