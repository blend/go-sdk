/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package r2trace

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	opentracingExt "github.com/opentracing/opentracing-go/ext"

	"github.com/opentracing/opentracing-go/mocktracer"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/r2"
	"github.com/blend/go-sdk/tracing"
)

func TestStart(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	reqTracer := Tracer(mockTracer)

	req := r2.New("https://foo.com/bar", r2.OptHeader(make(http.Header)))
	rtf := reqTracer.Start(req.Request)

	spanContext, err := mockTracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Request.Header))
	assert.Nil(err)
	mockSpanContext := spanContext.(mocktracer.MockSpanContext)

	span := rtf.(r2TraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(mockSpanContext.SpanID, mockSpan.SpanContext.SpanID)
	assert.Equal(tracing.OperationHTTPRequestOutgoing, mockSpan.OperationName)

	assert.Len(mockSpan.Tags(), 6)
	assert.Equal(tracing.SpanTypeHTTP, mockSpan.Tags()[tracing.TagKeySpanType])
	assert.Equal(opentracingExt.SpanKindRPCClientEnum, mockSpan.Tags()[string(opentracingExt.SpanKind)])
	assert.True(mockSpan.FinishTime.IsZero())
}

func TestStartNoHeader(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	reqTracer := Tracer(mockTracer)

	req := r2.New("https://foo.com/bar")
	rtf := reqTracer.Start(req.Request)

	spanContext, err := mockTracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Request.Header))
	assert.Nil(err)
	mockSpanContext := spanContext.(mocktracer.MockSpanContext)

	span := rtf.(r2TraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(mockSpanContext.SpanID, mockSpan.SpanContext.SpanID)
	assert.Equal(tracing.OperationHTTPRequestOutgoing, mockSpan.OperationName)

	assert.Len(mockSpan.Tags(), 6)
	assert.Equal(tracing.SpanTypeHTTP, mockSpan.Tags()[tracing.TagKeySpanType])
	assert.Equal(opentracingExt.SpanKindRPCClientEnum, mockSpan.Tags()[string(opentracingExt.SpanKind)])
	assert.True(mockSpan.FinishTime.IsZero())
}

func TestStartWithParentSpan(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	reqTracer := Tracer(mockTracer)

	parentSpan := mockTracer.StartSpan("test_op")
	ctx := opentracing.ContextWithSpan(context.Background(), parentSpan)

	req := r2.New("https://foo.com/bar", r2.OptContext(ctx))
	rtf := reqTracer.Start(req.Request)

	span := rtf.(r2TraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(tracing.OperationHTTPRequestOutgoing, mockSpan.OperationName)

	mockParentSpan := parentSpan.(*mocktracer.MockSpan)
	assert.Equal(mockSpan.ParentID, mockParentSpan.SpanContext.SpanID)
}

func TestStartParameterizedPath(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	reqTracer := Tracer(mockTracer)

	req := r2.New("https://foo.com/", r2.OptHeader(make(http.Header)), r2.OptPathParameterized("bar/:bar_id", map[string]string{"bar_id": "123"}))
	rtf := reqTracer.Start(req.Request)

	spanContext, err := mockTracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Request.Header))
	assert.Nil(err)
	mockSpanContext := spanContext.(mocktracer.MockSpanContext)

	span := rtf.(r2TraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(mockSpanContext.SpanID, mockSpan.SpanContext.SpanID)
	assert.Equal(tracing.OperationHTTPRequestOutgoing, mockSpan.OperationName)

	assert.Len(mockSpan.Tags(), 6)
	assert.Equal(tracing.SpanTypeHTTP, mockSpan.Tags()[tracing.TagKeySpanType])
	assert.Equal(opentracingExt.SpanKindRPCClientEnum, mockSpan.Tags()[string(opentracingExt.SpanKind)])
	assert.Equal("https://foo.com/bar/:bar_id", mockSpan.Tags()[tracing.TagKeyResourceName])
	assert.True(mockSpan.FinishTime.IsZero())
}

func TestStartParameterizedPathAndServiceHostName(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	reqTracer := Tracer(mockTracer)

	req := r2.New("https://foo.com/", r2.OptHeader(make(http.Header)), r2.OptPathParameterized("bar/:bar_id", map[string]string{"bar_id": "123"}), r2.OptServiceHostName("somehost"))
	rtf := reqTracer.Start(req.Request)

	spanContext, err := mockTracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Request.Header))
	assert.Nil(err)
	mockSpanContext := spanContext.(mocktracer.MockSpanContext)

	span := rtf.(r2TraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(mockSpanContext.SpanID, mockSpan.SpanContext.SpanID)
	assert.Equal(tracing.OperationHTTPRequestOutgoing, mockSpan.OperationName)

	assert.Len(mockSpan.Tags(), 6)
	assert.Equal(tracing.SpanTypeHTTP, mockSpan.Tags()[tracing.TagKeySpanType])
	assert.Equal(opentracingExt.SpanKindRPCClientEnum, mockSpan.Tags()[string(opentracingExt.SpanKind)])
	assert.Equal("https://{somehost}/bar/:bar_id", mockSpan.Tags()[tracing.TagKeyResourceName])
	assert.True(mockSpan.FinishTime.IsZero())
}

func TestStartServiceHostName(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	reqTracer := Tracer(mockTracer)

	req := r2.New("https://foo.com/og/url", r2.OptHeader(make(http.Header)), r2.OptServiceHostName("somehost"))
	rtf := reqTracer.Start(req.Request)

	spanContext, err := mockTracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Request.Header))
	assert.Nil(err)
	mockSpanContext := spanContext.(mocktracer.MockSpanContext)

	span := rtf.(r2TraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(mockSpanContext.SpanID, mockSpan.SpanContext.SpanID)
	assert.Equal(tracing.OperationHTTPRequestOutgoing, mockSpan.OperationName)

	assert.Len(mockSpan.Tags(), 6)
	assert.Equal(tracing.SpanTypeHTTP, mockSpan.Tags()[tracing.TagKeySpanType])
	assert.Equal(opentracingExt.SpanKindRPCClientEnum, mockSpan.Tags()[string(opentracingExt.SpanKind)])
	assert.Equal("https://{somehost}/og/url", mockSpan.Tags()[tracing.TagKeyResourceName])
	assert.True(mockSpan.FinishTime.IsZero())
}

func TestFinish(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	reqTracer := Tracer(mockTracer)

	req := r2.New("https://foo.com/bar")
	rtf := reqTracer.Start(req.Request)
	rtf.Finish(req.Request, &http.Response{StatusCode: 200}, time.Now(), nil)

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
	rtf := reqTracer.Start(req.Request)
	rtf.Finish(req.Request, &http.Response{StatusCode: 500}, time.Now(), fmt.Errorf("error"))

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
