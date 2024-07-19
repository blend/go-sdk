/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package redistrace

import (
	"context"
	"strings"
	"time"

	opentracing "github.com/opentracing/opentracing-go"

	"github.com/blend/go-sdk/redis"
	"github.com/blend/go-sdk/tracing"
)

var (
	_ redis.Tracer = (*redisTracer)(nil)
)

// Tracer returns a new redis tracer.
func Tracer(tracer opentracing.Tracer) redis.Tracer {
	return redisTracer{tracer: tracer}
}

type redisTracer struct {
	tracer opentracing.Tracer
}

func (rt redisTracer) Do(ctx context.Context, cfg redis.Config, op string, args []string) redis.TraceFinisher {
	host, port := splitHostPort(cfg.Addr)
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeyResourceName, Value: op},
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeRedis},
		opentracing.Tag{Key: tracing.TagKeyTargetHost, Value: host},
		opentracing.Tag{Key: tracing.TagKeyTargetPort, Value: port},
		opentracing.Tag{Key: "out.db", Value: cfg.DB},
		tracing.TagMeasured(),
		opentracing.StartTime(time.Now().UTC()),
	}
	span, _ := tracing.StartSpanFromContext(ctx, rt.tracer, tracing.OperationRedisCommand, startOptions...)
	return redisTraceFinisher{span: span}
}

type redisTraceFinisher struct {
	span opentracing.Span
}

func (rtf redisTraceFinisher) Finish(ctx context.Context, err error) {
	if rtf.span == nil {
		return
	}
	tracing.SpanError(rtf.span, err)
	rtf.span.Finish()
}

func splitHostPort(addr string) (host, port string) {
	parts := strings.SplitN(addr, ":", 2)
	if len(parts) > 0 {
		host = parts[0]
	}
	if len(parts) > 1 {
		port = parts[1]
	}
	return
}
