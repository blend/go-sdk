package grpctrace

import (
	"context"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/blend/go-sdk/grpcutil"
	"github.com/blend/go-sdk/stats/tracing"
)

// TracedServerUnary returns a unary server interceptor.
func TracedServerUnary(tracer opentracing.Tracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, args interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if tracer == nil {
			return handler(ctx, args)
		}

		startTime := time.Now().UTC()
		md, ok := metadata.FromIncomingContext(ctx)
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
			opentracing.Tag{Key: tracing.TagKeyGRPCCallingConvention, Value: "unary"},
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
		ctx = opentracing.ContextWithSpan(ctx, span)

		var err error
		defer func() {
			if err != nil {
				tracing.SpanError(span, err)
			}
			span.Finish()
		}()
		var result interface{}
		result, err = handler(ctx, args)
		return result, err
	}
}
