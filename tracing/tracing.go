package tracing

import (
	"context"
	"fmt"

	opentracing "github.com/opentracing/opentracing-go"

	"github.com/blend/go-sdk/ex"
)

// StartSpanFromContext creates a new span from a given context.
// It is required because opentracing relies on global state.
func StartSpanFromContext(ctx context.Context, tracer opentracing.Tracer, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	}
	span := tracer.StartSpan(operationName, opts...)
	return span, opentracing.ContextWithSpan(ctx, span)
}

// GetTracingSpanFromContext returns a tracing span from a given context.
func GetTracingSpanFromContext(ctx context.Context, key interface{}) opentracing.Span {
	if typed, ok := ctx.Value(key).(opentracing.Span); ok {
		return typed
	}
	return nil
}

// SpanError injects error metadata into a span.
func SpanError(span opentracing.Span, err error) {
	if err != nil {
		if typed := ex.As(err); typed != nil {
			span.SetTag(TagKeyError, typed.Class)
			span.SetTag(TagKeyErrorMessage, typed.Message)
			span.SetTag(TagKeyErrorStack, typed.StackTrace.String())
		} else {
			span.SetTag(TagKeyError, fmt.Sprintf("%v", err))
		}
	}
}
