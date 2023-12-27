/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package oauthtrace

import (
	"context"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	opentracingExt "github.com/opentracing/opentracing-go/ext"

	"golang.org/x/oauth2"

	"github.com/blend/go-sdk/oauth"
	"github.com/blend/go-sdk/tracing"
)

var (
	_	oauth.Tracer		= (*oauthTracer)(nil)
	_	oauth.TraceFinisher	= (*oauthTraceFinisher)(nil)
)

// Tracer returns a request tracer that also injects span context into outgoing headers.
func Tracer(tracer opentracing.Tracer) oauth.Tracer {
	return &oauthTracer{tracer: tracer}
}

type oauthTracer struct {
	tracer opentracing.Tracer
}

func (t oauthTracer) Start(ctx context.Context, config *oauth2.Config) oauth.TraceFinisher {
	startOptions := []opentracing.StartSpanOption{
		opentracingExt.SpanKindRPCClient,
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeHTTP},
		tracing.TagMeasured(),
		opentracing.StartTime(time.Now().UTC()),
	}
	span, _ := tracing.StartSpanFromContext(ctx, t.tracer, tracing.OperationHTTPRequestOutgoing, startOptions...)
	return oauthTraceFinisher{span: span}
}

type oauthTraceFinisher struct {
	span opentracing.Span
}

func (of oauthTraceFinisher) Finish(ctx context.Context, config *oauth2.Config, result *oauth.Result, err error) {
	if of.span == nil {
		return
	}
	tracing.SpanError(of.span, err)
	if result != nil {
		of.span.SetTag(tracing.TagKeyOAuthUsername, result.Profile.Email)
	}
	of.span.Finish()
}
