package grpctrace

import (
	"context"
	"time"

	"github.com/blend/go-sdk/grpcutil"
	"github.com/blend/go-sdk/stats/tracing"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// TracedServerStream returns a grpc streaming interceptor.
func TracedServerStream(tracer opentracing.Tracer) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if tracer == nil {
			return handler(srv, ss)
		}

		startTime := time.Now().UTC()
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			md = metadata.New(nil)
		}

		authority := grpcutil.MetaValue(md, grpcutil.MetaTagAuthority)
		contentType := grpcutil.MetaValue(md, grpcutil.MetaTagContentType)
		userAgent := grpcutil.MetaValue(md, grpcutil.MetaTagUserAgent)

		startOptions := []opentracing.StartSpanOption{
			opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeGRPC},
			opentracing.Tag{Key: tracing.TagKeyResourceName, Value: info.FullMethod},
			opentracing.Tag{Key: tracing.TagKeyGRPCMethod, Value: info.FullMethod},
			opentracing.Tag{Key: tracing.TagKeyGRPCAuthority, Value: authority},
			opentracing.Tag{Key: tracing.TagKeyGRPCUserAgent, Value: userAgent},
			opentracing.Tag{Key: tracing.TagKeyGRPCContentType, Value: contentType},
			opentracing.Tag{Key: tracing.TagKeyGRPCRole, Value: "server"},
			opentracing.Tag{Key: tracing.TagKeyGRPCCallingConvention, Value: "stream"},
			opentracing.StartTime(startTime),
		}

		// try to extract an incoming span context
		// this is typically done if we're a service being called in a chain from another (more ancestral)
		// span context.
		spanContext, _ := tracer.Extract(opentracing.HTTPHeaders, MetadataReaderWriter{md})
		if spanContext != nil {
			startOptions = append(startOptions, opentracing.ChildOf(spanContext))
		}

		span := tracer.StartSpan(tracing.OperationRPC, startOptions...)
		var err error
		defer func() {
			if err != nil {
				tracing.SpanError(span, err)
			}
			span.Finish()
		}()
		err = handler(srv, &spanServerStream{ServerStream: ss, Span: span})
		return nil
	}
}

// spanServerStream wraps around the embedded grpc.ServerStream, and
// intercepts calls to `Context()` returning a context with the span information injected.
//
// NOTE: you can extend this type to intercept calls to `SendMsg` and `RecvMsg` if you want to
// add tracing handling for individual stream calls.
type spanServerStream struct {
	grpc.ServerStream
	Span opentracing.Span
}

func (ss *spanServerStream) Context() context.Context {
	return opentracing.ContextWithSpan(ss.ServerStream.Context(), ss.Span)
}
