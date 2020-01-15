package httptrace_test

import (
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/stats/tracing"
	"github.com/blend/go-sdk/stats/tracing/httptrace"
	"github.com/blend/go-sdk/webutil"
)

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
		tracing.TagKeyResourceName: path,
		tracing.TagKeySpanType:     tracing.SpanTypeWeb,
		tracing.TagKeyHTTPMethod:   "GET",
		tracing.TagKeyHTTPURL:      path,
		"http.remote_addr":         "127.0.0.1",
		"http.host":                "localhost",
		"http.user_agent":          "go-sdk test",
	}
	assert.Equal(expectedTags, mockSpan.Tags())
	assert.True(mockSpan.FinishTime.IsZero())
}
