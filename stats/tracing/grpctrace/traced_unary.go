package grpctrace

import (
	"context"
	"strings"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/blend/go-sdk/grpcutil"
	"github.com/blend/go-sdk/stats/tracing"
)

// TracedUnary returns a unary server interceptor.
func TracedUnary(tracer opentracing.Tracer) grpc.UnaryServerInterceptor {
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
			opentracing.StartTime(startTime),
		}

		// try to extract an incoming span context
		// this is typically done if we're a service being called in a chain from another (more ancestral)
		// span context.
		spanContext, _ := tracer.Extract(opentracing.HTTPHeaders, metadataReaderWriter{md})
		if spanContext != nil {
			startOptions = append(startOptions, opentracing.ChildOf(spanContext))
		}

		span, ctx := tracing.StartSpanFromContext(ctx, tracer, tracing.OperationRPC, startOptions...)
		defer span.Finish()

		result, err := handler(ctx, args)
		if err != nil {
			tracing.SpanError(span, err)
		}
		return result, err
	}
}

// metadataReaderWriter satisfies both the opentracing.TextMapReader and
// opentracing.TextMapWriter interfaces.
type metadataReaderWriter struct {
	metadata.MD
}

func (w metadataReaderWriter) Set(key, val string) {
	// The GRPC HPACK implementation rejects any uppercase keys here.
	//
	// As such, since the HTTP_HEADERS format is case-insensitive anyway, we
	// blindly lowercase the key (which is guaranteed to work in the
	// Inject/Extract sense per the OpenTracing spec).
	key = strings.ToLower(key)
	w.MD[key] = append(w.MD[key], val)
}

func (w metadataReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for k, vals := range w.MD {
		for _, v := range vals {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}

	return nil
}
