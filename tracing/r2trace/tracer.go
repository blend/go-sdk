/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package r2trace

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	opentracingExt "github.com/opentracing/opentracing-go/ext"

	"github.com/blend/go-sdk/r2"
	"github.com/blend/go-sdk/tracing"
)

var (
	_	r2.Tracer		= (*r2Tracer)(nil)
	_	r2.TraceFinisher	= (*r2TraceFinisher)(nil)
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
		opentracingExt.SpanKindRPCClient,
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeHTTP},
		opentracing.Tag{Key: tracing.TagKeyResourceName, Value: r2.GetRawURLParameterized(req)},
		opentracing.Tag{Key: tracing.TagKeyHTTPMethod, Value: strings.ToUpper(req.Method)},
		opentracing.Tag{Key: tracing.TagKeyHTTPURL, Value: req.URL.String()},
		tracing.TagMeasured(),
		opentracing.StartTime(time.Now().UTC()),
	}
	span, ctx := tracing.StartSpanFromContext(req.Context(), rt.tracer, tracing.OperationHTTPRequestOutgoing, startOptions...)
	*req = *req.WithContext(ctx)

	if req.Header == nil {
		req.Header = make(http.Header)
	}
	_ = rt.tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
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
	if res != nil {
		rtf.span.SetTag(tracing.TagKeyHTTPCode, strconv.Itoa(res.StatusCode))
	} else {
		rtf.span.SetTag(tracing.TagKeyHTTPCode, http.StatusInternalServerError)
	}
	rtf.span.Finish()
}
