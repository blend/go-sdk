/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webtrace

import (
	"fmt"
	"strconv"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	opentracingExt "github.com/opentracing/opentracing-go/ext"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/tracing"
	"github.com/blend/go-sdk/tracing/httptrace"
	"github.com/blend/go-sdk/web"
)

var (
	_	web.Tracer		= (*webTracer)(nil)
	_	web.TraceFinisher	= (*webTraceFinisher)(nil)
	_	web.ViewTracer		= (*webTracer)(nil)
	_	web.ViewTraceFinisher	= (*webViewTraceFinisher)(nil)
)

// TracerOption is a tracer option.
type TracerOption func(*webTracer)

// OptIncludeCtxLabels toggles the option to include ctx labels in DD APM trace spans.
func OptIncludeCtxLabels(includeCtxLabels bool) TracerOption {
	return func(tracer *webTracer) {
		tracer.logLabelsToCtx = includeCtxLabels
	}
}

// Tracer returns a web tracer.
func Tracer(tracer opentracing.Tracer, options ...TracerOption) web.Tracer {
	wt := &webTracer{tracer: tracer, logLabelsToCtx: false}
	for _, opt := range options {
		opt(wt)
	}
	return wt
}

type webTracer struct {
	logLabelsToCtx	bool
	tracer		opentracing.Tracer
}

func (wt webTracer) Start(ctx *web.Ctx) web.TraceFinisher {
	var resource string
	var extra []opentracing.StartSpanOption
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
	return &webTraceFinisher{span: span, logLabelsToCtx: wt.logLabelsToCtx}
}

type webTraceFinisher struct {
	logLabelsToCtx	bool
	span		opentracing.Span
}

func (wtf webTraceFinisher) Finish(ctx *web.Ctx, err error) {
	if wtf.span == nil {
		return
	}
	tracing.SpanError(wtf.span, err)
	wtf.span.SetTag(tracing.TagKeyHTTPCode, strconv.Itoa(ctx.Response.StatusCode()))
	// Initially defaulted to false to allow opt-in by initial apps
	// to bake in the new behavior.
	if wtf.logLabelsToCtx {
		// Add explicit log labels as context for DD trace spans.
		for label, labelValue := range logger.GetLabels(ctx.Context()) {
			wtf.span.SetTag(fmt.Sprintf("%s.%s", tracing.TagKeyCtx, label), labelValue)
		}
	}
	wtf.span.Finish()
}

func (wt webTracer) StartView(ctx *web.Ctx, vr *web.ViewResult) web.ViewTraceFinisher {
	// set up basic start options (these are mostly tags).
	startOptions := []opentracing.StartSpanOption{
		opentracingExt.SpanKindRPCServer,
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
