package grpcutil

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/blend/go-sdk/logger"
)

// LoggedServerUnary returns a unary server interceptor.
func LoggedServerUnary(log logger.Triggerable) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, args interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now().UTC()
		result, err := handler(ctx, args)
		if log != nil {
			event := NewRPCEvent(info.FullMethod, time.Now().UTC().Sub(startTime))
			event.Err = err
			if md, ok := metadata.FromIncomingContext(ctx); ok {
				event.Authority = MetaValue(md, MetaTagAuthority)
				event.UserAgent = MetaValue(md, MetaTagUserAgent)
				event.ContentType = MetaValue(md, MetaTagContentType)
			}
			log.TriggerContext(ctx, event)
		}
		return result, err
	}
}

// LoggedClientUnary returns a unary client interceptor.
func LoggedClientUnary(log logger.Triggerable) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		startTime := time.Now().UTC()
		err := invoker(ctx, method, req, reply, cc, opts...)
		if log != nil {
			event := NewRPCEvent(method, time.Now().UTC().Sub(startTime))
			event.Err = err
			if md, ok := metadata.FromOutgoingContext(ctx); ok {
				event.Authority = MetaValue(md, MetaTagAuthority)
				event.UserAgent = MetaValue(md, MetaTagUserAgent)
				event.ContentType = MetaValue(md, MetaTagContentType)
			}
			log.TriggerContext(ctx, event)
		}
		return err
	}
}

// LoggedServerStream returns a stream server interceptor.
func LoggedServerStream(log logger.Triggerable) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		startTime := time.Now().UTC()
		err = handler(srv, stream)
		if log != nil {
			event := NewRPCEvent(info.FullMethod, time.Now().UTC().Sub(startTime))
			event.Err = err
			if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
				event.Authority = MetaValue(md, MetaTagAuthority)
				event.UserAgent = MetaValue(md, MetaTagUserAgent)
				event.ContentType = MetaValue(md, MetaTagContentType)
			}
			log.TriggerContext(stream.Context(), event)
		}
		return err
	}
}

// LoggedClientStream returns a stream server interceptor.
func LoggedClientStream(log logger.Triggerable) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		startTime := time.Now().UTC()
		clientStreamer, err := streamer(ctx, desc, cc, method, opts...)
		if log != nil {
			event := NewRPCEvent(method, time.Now().UTC().Sub(startTime))
			event.Err = err
			if md, ok := metadata.FromOutgoingContext(ctx); ok {
				event.Authority = MetaValue(md, MetaTagAuthority)
				event.UserAgent = MetaValue(md, MetaTagUserAgent)
				event.ContentType = MetaValue(md, MetaTagContentType)
			}
			log.TriggerContext(ctx, event)
		}
		return clientStreamer, err
	}
}
