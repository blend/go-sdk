package r2trace

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/r2"
	"github.com/blend/go-sdk/stats/tracing"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
)

func TestStart(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	reqTracer := Tracer(mockTracer)

	req := r2.New("https://foo.com/bar", r2.OptHeader(make(http.Header)))
	rtf := reqTracer.Start(&req.Request)

	spanContext, err := mockTracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	assert.Nil(err)
	mockSpanContext := spanContext.(mocktracer.MockSpanContext)

	span := rtf.(r2TraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(mockSpanContext.SpanID, mockSpan.SpanContext.SpanID)
	assert.Equal(tracing.OperationHTTPRequest, mockSpan.OperationName)

	assert.Len(mockSpan.Tags(), 1)
	assert.Equal(tracing.SpanTypeHTTP, mockSpan.Tags()[tracing.TagKeySpanType])
	assert.True(mockSpan.FinishTime.IsZero())
}

func TestStartNoHeader(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	reqTracer := Tracer(mockTracer)

	req := r2.New("https://foo.com/bar")
	rtf := reqTracer.Start(&req.Request)

	spanContext, err := mockTracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	assert.Nil(err)
	mockSpanContext := spanContext.(mocktracer.MockSpanContext)

	span := rtf.(r2TraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(mockSpanContext.SpanID, mockSpan.SpanContext.SpanID)
	assert.Equal(tracing.OperationHTTPRequest, mockSpan.OperationName)

	assert.Len(mockSpan.Tags(), 1)
	assert.Equal(tracing.SpanTypeHTTP, mockSpan.Tags()[tracing.TagKeySpanType])
	assert.True(mockSpan.FinishTime.IsZero())
}

func TestStartWithParentSpan(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	reqTracer := Tracer(mockTracer)

	parentSpan := mockTracer.StartSpan("test_op")
	ctx := opentracing.ContextWithSpan(context.Background(), parentSpan)

	req := r2.New("https://foo.com/bar", r2.OptContext(ctx))
	rtf := reqTracer.Start(&req.Request)

	span := rtf.(r2TraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(tracing.OperationHTTPRequest, mockSpan.OperationName)

	mockParentSpan := parentSpan.(*mocktracer.MockSpan)
	assert.Equal(mockSpan.ParentID, mockParentSpan.SpanContext.SpanID)
}

func TestFinish(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	reqTracer := Tracer(mockTracer)

	req := r2.New("https://foo.com/bar")
	rtf := reqTracer.Start(&req.Request)
	rtf.Finish(&req.Request, &http.Response{StatusCode: 200}, time.Now(), nil)

	span := rtf.(r2TraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal("200", mockSpan.Tags()[tracing.TagKeyHTTPCode])
	assert.False(mockSpan.FinishTime.IsZero())
}

func TestFinishError(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	reqTracer := Tracer(mockTracer)

	req := r2.New("https://foo.com/bar")
	rtf := reqTracer.Start(&req.Request)
	rtf.Finish(&req.Request, &http.Response{StatusCode: 500}, time.Now(), fmt.Errorf("error"))

	span := rtf.(r2TraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal("500", mockSpan.Tags()[tracing.TagKeyHTTPCode])
	assert.Equal("error", mockSpan.Tags()[tracing.TagKeyError])
	assert.False(mockSpan.FinishTime.IsZero())
}

func TestFinishNilSpan(t *testing.T) {
	assert := assert.New(t)

	rtf := r2TraceFinisher{}
	rtf.Finish(nil, nil, time.Now(), nil)
	assert.Nil(rtf.span)
}
