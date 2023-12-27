/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package tracing

import (
	"context"
	"strconv"

	"github.com/blend/go-sdk/logger"

	opentracing "github.com/opentracing/opentracing-go"
)

type tracerKey struct{}

// WithTracer adds a tracer to a context.
func WithTracer(ctx context.Context, tracer opentracing.Tracer) context.Context {
	return context.WithValue(ctx, tracerKey{}, tracer)
}

// GetTracer gets a tracer from a context.
func GetTracer(ctx context.Context) opentracing.Tracer {
	if value := ctx.Value(tracerKey{}); value != nil {
		if typed, ok := value.(opentracing.Tracer); ok {
			return typed
		}
	}
	return nil
}

// WithTraceAnnotations extracts trace span details as logger annotations onto a context
func WithTraceAnnotations(ctx context.Context, span opentracing.SpanContext) context.Context {
	if spanIDProvider, ok := span.(SpanIDProvider); ok {
		ctx = logger.WithAnnotation(ctx, LoggerAnnotationTracingSpanID, strconv.FormatUint(spanIDProvider.SpanID(), 10))
	}
	if traceIDProvider, ok := span.(TraceIDProvider); ok {
		ctx = logger.WithAnnotation(ctx, LoggerAnnotationTracingTraceID, strconv.FormatUint(traceIDProvider.TraceID(), 10))
	}
	return ctx
}
