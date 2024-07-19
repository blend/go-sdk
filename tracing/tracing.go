/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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
	ctx = WithTraceAnnotations(ctx, span.Context())
	return span, opentracing.ContextWithSpan(ctx, span)
}

// GetTracingSpanFromContext returns a tracing span from a given context.
func GetTracingSpanFromContext(ctx context.Context) opentracing.Span {
	return opentracing.SpanFromContext(ctx)
}

// Background returns a new `context.Background()`
// with the parent span from a given context.
//
// It is useful if you want to kick out goroutines but
// maintain tracing data.
func Background(ctx context.Context) context.Context {
	output := context.Background()
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		return opentracing.ContextWithSpan(output, parentSpan)
	}
	return output
}

// SpanError injects error metadata into a span.
func SpanError(span opentracing.Span, err error) {
	if err != nil {
		if typed := ex.As(err); typed != nil {
			span.SetTag(TagKeyError, typed.Class)
			span.SetTag(TagKeyErrorType, typed.Class)
			span.SetTag(TagKeyErrorMessage, typed.Message)
			span.SetTag(TagKeyErrorStack, typed.StackTrace.String())
			span.SetTag(TagKeyErrorDetails, fmt.Sprintf("%+v", err))
		} else {
			span.SetTag(TagKeyError, fmt.Sprintf("%v", err))
		}
	}
}

// TagMeasured returns an opentracing tag to indicate the span should be measured.
func TagMeasured() opentracing.Tag {
	return opentracing.Tag{Key: TagKeyMeasured, Value: 1}
}
