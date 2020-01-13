package webtrace

import (
	"fmt"
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/stats/tracing"
	"github.com/blend/go-sdk/web"
)

func TestStart(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	webTracer := Tracer(mockTracer)

	ctx := web.MockCtx("GET", "/test-resource")
	_ = webTracer.Start(ctx)

	span := opentracing.SpanFromContext(ctx.Context())
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(tracing.OperationHTTPRequest, mockSpan.OperationName)

	assert.Len(mockSpan.Tags(), 7)
	assert.Equal("/test-resource", mockSpan.Tags()[tracing.TagKeyResourceName])
	assert.Equal(tracing.SpanTypeWeb, mockSpan.Tags()[tracing.TagKeySpanType])
	assert.Equal("GET", mockSpan.Tags()[tracing.TagKeyHTTPMethod])
	assert.Equal("/test-resource", mockSpan.Tags()[tracing.TagKeyHTTPURL])
	assert.Equal("127.0.0.1", mockSpan.Tags()["http.remote_addr"])
	assert.Equal("localhost", mockSpan.Tags()["http.host"])
	assert.Equal("go-sdk test", mockSpan.Tags()["http.user_agent"])
	assert.True(mockSpan.FinishTime.IsZero())
}

func TestStartWithRoute(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	webTracer := Tracer(mockTracer)

	ctx := web.MockCtx(
		"GET", "/test-resource/3",
		web.OptCtxRoute(&web.Route{Path: "/test-resource/:id"}),
	)
	_ = webTracer.Start(ctx)

	span := opentracing.SpanFromContext(ctx.Context())
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(tracing.OperationHTTPRequest, mockSpan.OperationName)

	assert.Len(mockSpan.Tags(), 8)
	assert.Equal("/test-resource/:id", mockSpan.Tags()[tracing.TagKeyResourceName])
	assert.Equal(tracing.SpanTypeWeb, mockSpan.Tags()[tracing.TagKeySpanType])
	assert.Equal("GET", mockSpan.Tags()[tracing.TagKeyHTTPMethod])
	assert.Equal("/test-resource/3", mockSpan.Tags()[tracing.TagKeyHTTPURL])
	assert.Equal("127.0.0.1", mockSpan.Tags()["http.remote_addr"])
	assert.Equal("localhost", mockSpan.Tags()["http.host"])
	assert.Equal("go-sdk test", mockSpan.Tags()["http.user_agent"])
	assert.Equal("/test-resource/:id", mockSpan.Tags()["http.route"])
	assert.True(mockSpan.FinishTime.IsZero())
}

func optCtxIncomingSpan(t opentracing.Tracer, s opentracing.Span) web.CtxOption {
	return func(c *web.Ctx) {
		t.Inject(
			s.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request.Header),
		)
	}
}

func TestStartWithParentSpan(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	webTracer := Tracer(mockTracer)

	parentSpan := mockTracer.StartSpan("test_op")
	ctx := web.MockCtx(
		"GET", "/test-resource",
		optCtxIncomingSpan(mockTracer, parentSpan),
	)
	_ = webTracer.Start(ctx)

	span := opentracing.SpanFromContext(ctx.Context())
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(tracing.OperationHTTPRequest, mockSpan.OperationName)

	mockParentSpan := parentSpan.(*mocktracer.MockSpan)
	assert.Equal(mockSpan.ParentID, mockParentSpan.SpanContext.SpanID)
}

func TestFinish(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	webTracer := Tracer(mockTracer)

	ctx := web.MockCtx("GET", "/test-resource")
	tf := webTracer.Start(ctx)

	result := web.Raw([]byte("success"))
	result.Render(ctx)

	tf.Finish(ctx, nil)

	span := opentracing.SpanFromContext(ctx.Context())
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal("200", mockSpan.Tags()[tracing.TagKeyHTTPCode])
	assert.False(mockSpan.FinishTime.IsZero())
}

func TestFinishError(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	webTracer := Tracer(mockTracer)

	ctx := web.MockCtx("GET", "/test-resource")
	tf := webTracer.Start(ctx)

	result := web.Raw([]byte("success"))
	result.StatusCode = 500
	result.Render(ctx)

	tf.Finish(ctx, fmt.Errorf("error"))

	span := opentracing.SpanFromContext(ctx.Context())
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal("500", mockSpan.Tags()[tracing.TagKeyHTTPCode])
	assert.Equal("error", mockSpan.Tags()[tracing.TagKeyError])
	assert.False(mockSpan.FinishTime.IsZero())
}

func TestFinishNilSpan(t *testing.T) {
	assert := assert.New(t)

	ctx := web.MockCtx("GET", "/test-resource")
	webTraceFinisher{}.Finish(ctx, nil)
	assert.Nil(opentracing.SpanFromContext(ctx.Context()))
}

func TestStartView(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	webViewTracer := Tracer(mockTracer).(web.ViewTracer)

	ctx := web.MockCtx("GET", "/test-resource")
	viewResult := &web.ViewResult{
		ViewName:   "test_view",
		StatusCode: 200,
	}
	wvtf := webViewTracer.StartView(ctx, viewResult)
	span := wvtf.(*webViewTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(tracing.OperationHTTPRender, mockSpan.OperationName)

	assert.Len(mockSpan.Tags(), 2)
	assert.Equal("test_view", mockSpan.Tags()[tracing.TagKeyResourceName])
	assert.Equal(tracing.SpanTypeWeb, mockSpan.Tags()[tracing.TagKeySpanType])
	assert.True(mockSpan.FinishTime.IsZero())
}

func TestStartViewWithParentSpan(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	webViewTracer := Tracer(mockTracer).(web.ViewTracer)

	parentSpan := mockTracer.StartSpan("test_op")
	ctx := web.MockCtx("GET", "/test-resource")
	ctx.WithContext(opentracing.ContextWithSpan(ctx.Context(), parentSpan))
	viewResult := &web.ViewResult{
		ViewName:   "test_view",
		StatusCode: 200,
	}
	wvtf := webViewTracer.StartView(ctx, viewResult)
	span := wvtf.(*webViewTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(tracing.OperationHTTPRender, mockSpan.OperationName)

	mockParentSpan := parentSpan.(*mocktracer.MockSpan)
	assert.Equal(mockSpan.ParentID, mockParentSpan.SpanContext.SpanID)
}

func TestFinishView(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	webViewTracer := Tracer(mockTracer).(web.ViewTracer)

	ctx := web.MockCtx("GET", "/test-resource")
	viewResult := &web.ViewResult{
		ViewName:   "test_view",
		StatusCode: 200,
	}
	wvtf := webViewTracer.StartView(ctx, viewResult)
	wvtf.FinishView(ctx, viewResult, nil)

	span := wvtf.(*webViewTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)

	assert.Nil(mockSpan.Tags()[tracing.TagKeyError])
	assert.False(mockSpan.FinishTime.IsZero())
}
