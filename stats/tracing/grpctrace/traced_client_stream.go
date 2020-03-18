package grpctrace

import (
	"context"
	"time"

	opentracing "github.com/opentracing/opentracing-go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/blend/go-sdk/stats/tracing"
)

// TracedClientStream returns a stream client interceptor that adds tracing spans.
func TracedClientStream(tracer opentracing.Tracer) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if tracer == nil {
			return streamer(ctx, desc, cc, method, opts...)
		}
		startTime := time.Now().UTC()

		startOptions := []opentracing.StartSpanOption{
			opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeGRPC},
			opentracing.Tag{Key: tracing.TagKeyResourceName, Value: method},
			opentracing.Tag{Key: tracing.TagKeyGRPCMethod, Value: method},
			opentracing.Tag{Key: tracing.TagKeyGRPCRole, Value: "client"},
			opentracing.Tag{Key: tracing.TagKeyGRPCCallingConvention, Value: "stream"},
			opentracing.StartTime(startTime),
		}
		span, ctx := tracing.StartSpanFromContext(ctx, tracer, tracing.OperationRPC, startOptions...)
		var err error
		var cs grpc.ClientStream
		defer func() {
			if err != nil {
				tracing.SpanError(span, err)
			}
			span.Finish()
		}()

		md := make(metadata.MD)
		tracer.Inject(span.Context(), opentracing.TextMap, MetadataReaderWriter{md})
		ctx = metadata.NewOutgoingContext(ctx, md)
		cs, err = streamer(ctx, desc, cc, method, opts...)
		return cs, err
	}
}
