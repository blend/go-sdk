/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package httptrace_test

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
	"github.com/blend/go-sdk/tracing"
	"github.com/blend/go-sdk/tracing/httptrace"
	"github.com/blend/go-sdk/webutil"
)

func TestStartHTTPSpan(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()

	path := "/test-resource"
	req := webutil.NewMockRequest("GET", path)
	resource := "/:id"
	startTime := time.Now().Add(-10 * time.Second)
	span, _ := httptrace.StartHTTPSpan(
		context.TODO(),
		mockTracer,
		req,
		resource,
		startTime,
		opentracing.Tag{Key: "http.route", Value: resource},
	)

	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(tracing.OperationHTTPRequest, mockSpan.OperationName)

	expectedTags := map[string]interface{}{
		tracing.TagKeyMeasured:          1,
		tracing.TagKeyResourceName:      fmt.Sprintf("GET %s", resource),
		tracing.TagKeySpanType:          tracing.SpanTypeWeb,
		tracing.TagKeyHTTPMethod:        "GET",
		tracing.TagKeyHTTPURL:           path,
		"http.remote_addr":              "127.0.0.1",
		"http.host":                     "localhost",
		"http.user_agent":               "go-sdk test",
		"http.route":                    resource,
		string(opentracingExt.SpanKind): opentracingExt.SpanKindRPCServerEnum,
	}
	assert.Equal(expectedTags, mockSpan.Tags())
	assert.Equal(startTime, mockSpan.StartTime)
	assert.True(mockSpan.FinishTime.IsZero())
}

func TestStart(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	httpTracer := httptrace.Tracer(mockTracer)

	path := "/test-resource"
	req := webutil.NewMockRequest("GET", path)
	_, req = httpTracer.Start(req)

	span := opentracing.SpanFromContext(req.Context())
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(tracing.OperationHTTPRequest, mockSpan.OperationName)

	expectedTags := map[string]interface{}{
		tracing.TagKeyMeasured:          1,
		tracing.TagKeyResourceName:      fmt.Sprintf("GET %s", path),
		tracing.TagKeySpanType:          tracing.SpanTypeWeb,
		tracing.TagKeyHTTPMethod:        "GET",
		tracing.TagKeyHTTPURL:           path,
		"http.remote_addr":              "127.0.0.1",
		"http.host":                     "localhost",
		"http.user_agent":               "go-sdk test",
		string(opentracingExt.SpanKind): opentracingExt.SpanKindRPCServerEnum,
	}
	assert.Equal(expectedTags, mockSpan.Tags())
	assert.True(mockSpan.FinishTime.IsZero())
}

func applyIncomingSpan(req *http.Request, t opentracing.Tracer, s opentracing.Span) {
	_ = t.Inject(
		s.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
}

func TestStartWithParentSpan(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	httpTracer := httptrace.Tracer(mockTracer)

	parentSpan := mockTracer.StartSpan("test_op")
	path := "/test-resource"
	req := webutil.NewMockRequest("GET", path)
	applyIncomingSpan(req, mockTracer, parentSpan)
	_, req = httpTracer.Start(req)

	span := opentracing.SpanFromContext(req.Context())
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(tracing.OperationHTTPRequest, mockSpan.OperationName)

	mockParentSpan := parentSpan.(*mocktracer.MockSpan)
	assert.Equal(mockSpan.ParentID, mockParentSpan.SpanContext.SpanID)
}

func TestFinish(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	httpTracer := httptrace.Tracer(mockTracer)

	path := "/test-resource"
	req := webutil.NewMockRequest("GET", path)
	tf, req := httpTracer.Start(req)

	tf.Finish(http.StatusOK, nil)

	span := opentracing.SpanFromContext(req.Context())
	mockSpan := span.(*mocktracer.MockSpan)
	assert.False(mockSpan.FinishTime.IsZero())
}

func TestFinishError(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	httpTracer := httptrace.Tracer(mockTracer)

	path := "/test-resource"
	req := webutil.NewMockRequest("GET", path)
	tf, req := httpTracer.Start(req)

	tf.Finish(http.StatusOK, fmt.Errorf("error"))

	span := opentracing.SpanFromContext(req.Context())
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal("error", mockSpan.Tags()[tracing.TagKeyError])
	assert.Equal("200", mockSpan.Tags()[tracing.TagKeyHTTPCode])
	assert.False(mockSpan.FinishTime.IsZero())
}
