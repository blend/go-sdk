package webtrace

import (
	"strconv"
	"time"
	"github.com/blend/go-sdk/stats/tracing"
	"github.com/blend/go-sdk/web"
	"github.com/blend/go-sdk/webutil"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	// StateKeySpan is the span state key.
	StateKeySpan = "web-span"
)

var (
	_ web.Tracer            = (*webTracer)(nil)
	_ web.TraceFinisher     = (*webTraceFinisher)(nil)
	_ web.ViewTracer        = (*webTracer)(nil)
	_ web.ViewTraceFinisher = (*webViewTraceFinisher)(nil)
)

// Tracer returns a web tracer.
func Tracer(tracer opentracing.Tracer) web.Tracer {
	return &webTracer{tracer: tracer}
}

type webTracer struct {
	tracer opentracing.Tracer
}

func (wt webTracer) Start(ctx *web.Ctx) web.TraceFinisher {
	var resource string
	if ctx.Route != nil {
		resource = ctx.Route.String()
	} else {
		resource = ctx.Request.URL.Path
	}

	// set up basic start options (these are mostly tags).
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeyResourceName, Value: resource},
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeWeb},
		opentracing.Tag{Key: tracing.TagKeyHTTPMethod, Value: ctx.Request.Method},
		opentracing.Tag{Key: tracing.TagKeyHTTPURL, Value: ctx.Request.URL.Path},
		opentracing.Tag{Key: "http.remote_addr", Value: webutil.GetRemoteAddr(ctx.Request)},
		opentracing.Tag{Key: "http.host", Value: webutil.GetHost(ctx.Request)},
		opentracing.Tag{Key: "http.user_agent", Value: webutil.GetUserAgent(ctx.Request)},
		opentracing.StartTime(ctx.RequestStart),
	}
	if ctx.Route != nil {
		startOptions = append(startOptions, opentracing.Tag{Key: "http.route", Value: ctx.Route.String()})
	}

	// try to extract an incoming span context
	// this is typically done if we're a service being called in a chain from another (more ancestral)
	// span context.
	spanContext, _ := wt.tracer.Extract(opentracing.TextMap, opentracing.HTTPHeadersCarrier(ctx.Request.Header))
	if spanContext != nil {
		startOptions = append(startOptions, opentracing.ChildOf(spanContext))
	}
	// start the span.
	span, spanCtx := tracing.StartSpanFromContext(ctx.Context(), wt.tracer, tracing.OperationHTTPRequest, startOptions...)
	// inject the new context
	ctx.Request = ctx.Request.WithContext(spanCtx)
	ctx.WithContext(spanCtx)
	return &webTraceFinisher{span: span}
}

type webTraceFinisher struct {
	span opentracing.Span
}

func (wtf webTraceFinisher) Finish(ctx *web.Ctx, err error) {
	if wtf.span == nil {
		return
	}
	tracing.SpanError(wtf.span, err)
	wtf.span.SetTag(tracing.TagKeyHTTPCode, strconv.Itoa(ctx.Response.StatusCode()))
	wtf.span.Finish()
}

func (wt webTracer) StartView(ctx *web.Ctx, vr *web.ViewResult) web.ViewTraceFinisher {
	// set up basic start options (these are mostly tags).
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeyResourceName, Value: vr.ViewName},
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeWeb},
		opentracing.StartTime(time.Now().UTC()),
	}
	// start the span.
	span, _ := tracing.StartSpanFromContext(ctx.Context(), wt.tracer, tracing.OperationHTTPRender, startOptions...)
	return &webViewTraceFinisher{span: span}
}

type webViewTraceFinisher struct {
	span opentracing.Span
}

func (wvtf webViewTraceFinisher) FinishView(ctx *web.Ctx, vr *web.ViewResult, err error) {
	if wvtf.span == nil {
		return
	}
	tracing.SpanError(wvtf.span, err)
	wvtf.span.Finish()
}
