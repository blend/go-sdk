package oauthtrace

import (
	"context"
	"testing"

	"github.com/opentracing/opentracing-go/mocktracer"
	"golang.org/x/oauth2"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/oauth"
	"github.com/blend/go-sdk/tracing"
)

func TestStart(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	oauthTracer := Tracer(mockTracer)

	finisher := oauthTracer.Start(context.Background(), &oauth2.Config{RedirectURL: "/admin"})
	span := finisher.(oauthTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(tracing.OperationHTTPRequest, mockSpan.OperationName)

	assert.Len(mockSpan.Tags(), 1)
	assert.Equal(tracing.SpanTypeHTTP, mockSpan.Tags()[tracing.TagKeySpanType])
	assert.True(mockSpan.FinishTime.IsZero())
}

func TestFinish(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	oauthTracer := Tracer(mockTracer)

	finisher := oauthTracer.Start(context.Background(), &oauth2.Config{RedirectURL: "/admin"})
	finisher.Finish(context.Background(), &oauth2.Config{RedirectURL: "/admin"}, &oauth.Result{Profile: oauth.Profile{Email: "example-string@blend.com"}}, nil)

	span := finisher.(oauthTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal("example-string@blend.com", mockSpan.Tags()[tracing.TagKeyOAuthUsername])
	assert.False(mockSpan.FinishTime.IsZero())
}
