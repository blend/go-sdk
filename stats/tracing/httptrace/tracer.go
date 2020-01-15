package httptrace

import (
	"net/http"
	"time"

	opentracing "github.com/opentracing/opentracing-go"

	"github.com/blend/go-sdk/stats/tracing"
	"github.com/blend/go-sdk/web"
	"github.com/blend/go-sdk/webutil"
)

var (
	_ web.HTTPTracer        = (*httpTracer)(nil)
	_ web.HTTPTraceFinisher = (*httpTraceFinisher)(nil)
)

// Tracer returns an HTTP tracer.
func Tracer(tracer opentracing.Tracer) web.HTTPTracer {
	return &httpTracer{tracer: tracer}
}

type httpTracer struct {
	tracer opentracing.Tracer
}

// Start opens a span and creates a new request with a modified context, based
// on the span that was opened.
func (ht httpTracer) Start(req *http.Request, resource string, startTime *time.Time) (web.HTTPTraceFinisher, *http.Request) {
	if resource == "" {
		resource = req.URL.Path
	}
	if startTime == nil {
		now := time.Now().UTC()
		startTime = &now
	}
	// set up basic start options (these are mostly tags).
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeyResourceName, Value: resource},
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeWeb},
		opentracing.Tag{Key: tracing.TagKeyHTTPMethod, Value: req.Method},
		opentracing.Tag{Key: tracing.TagKeyHTTPURL, Value: req.URL.Path},
		opentracing.Tag{Key: "http.remote_addr", Value: webutil.GetRemoteAddr(req)},
		opentracing.Tag{Key: "http.host", Value: webutil.GetHost(req)},
		opentracing.Tag{Key: "http.user_agent", Value: webutil.GetUserAgent(req)},
		opentracing.StartTime(*startTime),
	}

	// try to extract an incoming span context
	// this is typically done if we're a service being called in a chain from another (more ancestral)
	// span context.
	spanContext, _ := ht.tracer.Extract(opentracing.TextMap, opentracing.HTTPHeadersCarrier(req.Header))
	if spanContext != nil {
		startOptions = append(startOptions, opentracing.ChildOf(spanContext))
	}
	// start the span.
	span, spanCtx := tracing.StartSpanFromContext(req.Context(), ht.tracer, tracing.OperationHTTPRequest, startOptions...)
	// inject the new context
	newReq := req.WithContext(spanCtx)
	return &httpTraceFinisher{span: span}, newReq
}

type httpTraceFinisher struct {
	span opentracing.Span
}

func (htf httpTraceFinisher) Finish(req *http.Request, err error) {
	if htf.span == nil {
		return
	}
	tracing.SpanError(htf.span, err)
	htf.span.Finish()
}
