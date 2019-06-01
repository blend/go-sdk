package tracing

import (
	"context"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
)

func TestStartSpanFromContext(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()

	span, ctx := StartSpanFromContext(
		context.Background(),
		mockTracer,
		"test.operation.one",
		opentracing.Tags(map[string]interface{}{"k1": "v1"}),
	)

	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal("test.operation.one", mockSpan.OperationName)
	assert.Len(mockSpan.Tags(), 1)
	assert.Equal("v1", mockSpan.Tags()["k1"])

	spanFromCtx := opentracing.SpanFromContext(ctx)
	mockSpanFromCtx := spanFromCtx.(*mocktracer.MockSpan)
	assert.Equal(mockSpan.String(), mockSpanFromCtx.String())
}

func TestStartSpanFromContextWithParent(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()

	// Start Parent Span from Background Context
	parentSpan, ctx := StartSpanFromContext(
		context.Background(),
		mockTracer,
		"test.operation.one",
	)

	mockParentSpan := parentSpan.(*mocktracer.MockSpan)

	// Start Child Span from Context with Parent Span
	span, ctx := StartSpanFromContext(
		ctx,
		mockTracer,
		"test.operation.two",
	)

	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal("test.operation.two", mockSpan.OperationName)

	spanFromCtx := opentracing.SpanFromContext(ctx)
	mockSpanFromCtx := spanFromCtx.(*mocktracer.MockSpan)
	assert.Equal(mockSpan.String(), mockSpanFromCtx.String())

	assert.Equal(mockSpan.ParentID, mockParentSpan.SpanContext.SpanID)
}

func TestGetTracingSpanFromContext(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()

	spanKey := struct{}{}
	span := mockTracer.StartSpan("test.operation")
	ctx := context.Background()
	ctx = context.WithValue(ctx, spanKey, span)

	spanFromCtx := GetTracingSpanFromContext(ctx, spanKey)

	mockSpan := span.(*mocktracer.MockSpan)
	mockSpanFromCtx := spanFromCtx.(*mocktracer.MockSpan)
	assert.Equal(mockSpan.String(), mockSpanFromCtx.String())
}

func TestGetTracingSpanFromContextMiss(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()

	spanKey := struct{}{}
	span := mockTracer.StartSpan("test.operation")
	ctx := context.Background()
	ctx = context.WithValue(ctx, spanKey, span)

	spanFromCtx := GetTracingSpanFromContext(ctx, "wrongKey")
	assert.Nil(spanFromCtx)
}

func TestSpanError(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()

	// Non go-sdk/ex error
	nonExErr := fmt.Errorf("Test Error")
	span := mockTracer.StartSpan("test.operation")
	SpanError(span, nonExErr)

	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal("Test Error", mockSpan.Tags()[TagKeyError])
	assert.Nil(mockSpan.Tags()[TagKeyErrorMessage])
	assert.Nil(mockSpan.Tags()[TagKeyErrorStack])

	// go-sdk/ex error
	nonExErr = ex.New("Test Ex Error").
		WithMessage("Test Message")
	span = mockTracer.StartSpan("test.operation")
	SpanError(span, nonExErr)

	mockSpan = span.(*mocktracer.MockSpan)
	assert.Equal("Test Ex Error", mockSpan.Tags()[TagKeyError])
	assert.Equal("Test Message", mockSpan.Tags()[TagKeyErrorMessage])
	assert.NotNil(mockSpan.Tags()[TagKeyErrorStack])
	assert.Contains(mockSpan.Tags()[TagKeyErrorStack].(string), "tracing_test.go")
}
