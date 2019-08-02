package r2trace

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/blend/go-sdk/r2"
	"github.com/blend/go-sdk/stats/tracing"
	opentracing "github.com/opentracing/opentracing-go"
)

var (
	_ r2.Tracer        = (*r2Tracer)(nil)
	_ r2.TraceFinisher = (*r2TraceFinisher)(nil)
)

// Tracer returns a request tracer that also injects span context into outgoing headers.
func Tracer(tracer opentracing.Tracer) r2.Tracer {
	return &r2Tracer{tracer: tracer}
}

type r2Tracer struct {
	tracer opentracing.Tracer
}

func (rt r2Tracer) Start(req *http.Request) r2.TraceFinisher {
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeHTTP},
		opentracing.Tag{Key: tracing.TagKeyResourceName, Value: req.URL.String()},
		opentracing.Tag{Key: tracing.TagKeyHTTPMethod, Value: strings.ToUpper(req.Method)},
		opentracing.Tag{Key: tracing.TagKeyHTTPURL, Value: req.URL.String()},
		opentracing.StartTime(time.Now().UTC()),
	}
	span, _ := tracing.StartSpanFromContext(req.Context(), rt.tracer, tracing.OperationHTTPRequest, startOptions...)

	if req.Header == nil {
		req.Header = make(http.Header)
	}
	rt.tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	return r2TraceFinisher{span: span}
}

type r2TraceFinisher struct {
	span opentracing.Span
}

func (rtf r2TraceFinisher) Finish(req *http.Request, res *http.Response, ts time.Time, err error) {
	if rtf.span == nil {
		return
	}
	tracing.SpanError(rtf.span, err)
	rtf.span.SetTag(tracing.TagKeyHTTPCode, strconv.Itoa(res.StatusCode))
	rtf.span.Finish()
}
