package grpcutil

import (
	"context"
	"time"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"github.com/blend/go-sdk/logger"
)

// LoggedUnary returns a unary server interceptor.
func LoggedUnary(log logger.Triggerable) grpc.UnaryServerInterceptor {
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
			log.Trigger(ctx, event)
		}
		return result, err
	}
}

// LoggedStreaming returns a streaming server interceptor.
func LoggedStreaming(log logger.Triggerable) grpc.StreamServerInterceptor {
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
			log.Trigger(context.Background(), event)
		}
		return err
	}
}
