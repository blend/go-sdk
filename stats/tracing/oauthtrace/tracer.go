package oauthtrace

import (
	"net/http"
	"time"

	"github.com/blend/go-sdk/oauth"
	"github.com/blend/go-sdk/stats/tracing"
	opentracing "github.com/opentracing/opentracing-go"
	"golang.org/x/oauth2"
)

// Tracer returns a request tracer that also injects span context into outgoing headers.
func Tracer(tracer opentracing.Tracer) oauth.Tracer {
	return &oauthTracer{tracer: tracer}
}

type oauthTracer struct {
	tracer opentracing.Tracer
}

func (t oauthTracer) Start(config *oauth2.Config) oauth.TraceFinisher {
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeHTTP},
		opentracing.StartTime(time.Now().UTC()),
	}
	span, _ := tracing.StartSpanFromContext(req.Context(), t.tracer, tracing.OperationHTTPRequest, startOptions...)
	return oauthTraceFinisher{span: span}
}

type oauthTraceFinisher struct {
	span opentracing.Span
}

func (of oauthTraceFinisher) Finish(req *http.Request, result *oauth.Result, err error) {
	if of.span == nil {
		return
	}
	tracing.SpanError(of.span, err)
	of.span.SetTag(tracing.TagKeyOAuthUsername, result.Profile.Email)
	of.span.Finish()
}
