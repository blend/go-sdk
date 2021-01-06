package webtrace

import (
	"strconv"
	"time"

	opentracing "github.com/opentracing/opentracing-go"

	"github.com/blend/go-sdk/tracing"
	"github.com/blend/go-sdk/tracing/httptrace"
	"github.com/blend/go-sdk/web"
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
	extra := []opentracing.StartSpanOption{}
	if ctx.Route != nil {
		resource = ctx.Route.String()
		extra = append(extra, opentracing.Tag{Key: "http.route", Value: ctx.Route.String()})
	} else {
		resource = ctx.Request.URL.Path
	}
	span, newReq := httptrace.StartHTTPSpan(
		ctx.Context(),
		wt.tracer,
		ctx.Request,
		resource,
		ctx.RequestStarted,
		extra...,
	)
	ctx.Request = newReq
	ctx.WithContext(newReq.Context())
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
		tracing.TagMeasured(),
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
