package grpctrace

import (
	"context"
	"time"

	opentracing "github.com/opentracing/opentracing-go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/blend/go-sdk/stats/tracing"
)

// TracedClientUnary returns a unary client interceptor that adds tracing spans.
func TracedClientUnary(tracer opentracing.Tracer) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		if tracer == nil {
			err = invoker(ctx, method, req, reply, cc, opts...)
			return
		}
		startTime := time.Now().UTC()

		startOptions := []opentracing.StartSpanOption{
			opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeGRPC},
			opentracing.Tag{Key: tracing.TagKeyResourceName, Value: method},
			opentracing.Tag{Key: tracing.TagKeyGRPCMethod, Value: method},
			opentracing.Tag{Key: tracing.TagKeyGRPCRole, Value: "client"},

			opentracing.StartTime(startTime),
		}
		span, ctx := tracing.StartSpanFromContext(ctx, tracer, tracing.OperationRPC, startOptions...)
		defer func() {
			if err != nil {
				tracing.SpanError(span, err)
			}
			span.Finish()
		}()
		md := make(metadata.MD)
		tracer.Inject(span.Context(), opentracing.TextMap, MetadataReaderWriter{md})
		ctx = metadata.NewOutgoingContext(ctx, md)
		err = invoker(ctx, method, req, reply, cc, opts...)
		return
	}
}
