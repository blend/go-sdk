/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package tracing

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"

	opentracing "github.com/opentracing/opentracing-go"
)

type spanContextWithoutGetters struct {
	opentracing.SpanContext
}

type spanContextWithSpanID struct {
	opentracing.SpanContext
	spanID uint64
}

func (c spanContextWithSpanID) SpanID() uint64 {
	return c.spanID
}

type spanContextWithTraceID struct {
	opentracing.SpanContext
	traceID uint64
}

func (c spanContextWithTraceID) TraceID() uint64 {
	return c.traceID
}

type spanContextWithAllGetters struct {
	opentracing.SpanContext
	spanID  uint64
	traceID uint64
}

func (c spanContextWithAllGetters) SpanID() uint64 {
	return c.spanID
}

func (c spanContextWithAllGetters) TraceID() uint64 {
	return c.traceID
}

func TestWithTraceAnnotations_NoGetters(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	ctx = WithTraceAnnotations(ctx, spanContextWithoutGetters{})
	assert.Empty(logger.GetAnnotations(ctx))
}

func TestWithTraceAnnotations_AllGetters(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	ctx = WithTraceAnnotations(ctx, spanContextWithAllGetters{
		spanID:  123,
		traceID: 456,
	})

	annotations := logger.GetAnnotations(ctx)
	assert.Len(annotations, 2)
	assert.Equal("123", annotations[LoggerAnnotationTracingSpanID])
	assert.Equal("456", annotations[LoggerAnnotationTracingTraceID])
}

func TestWithTraceAnnotations_SpanIDProvider(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	ctx = WithTraceAnnotations(ctx, spanContextWithSpanID{
		spanID: 123,
	})

	annotations := logger.GetAnnotations(ctx)
	assert.Len(annotations, 1)
	assert.Equal("123", annotations[LoggerAnnotationTracingSpanID])
}

func TestWithTraceAnnotations_TraceIDProvider(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	ctx = WithTraceAnnotations(ctx, spanContextWithTraceID{
		traceID: 456,
	})

	annotations := logger.GetAnnotations(ctx)
	assert.Len(annotations, 1)
	assert.Equal("456", annotations[LoggerAnnotationTracingTraceID])
}
