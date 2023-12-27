/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package crontrace

import (
	"context"
	"time"

	opentracing "github.com/opentracing/opentracing-go"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/tracing"
)

// Tracer returns a opentracing cron tracer.
func Tracer(t opentracing.Tracer) cron.Tracer {
	return &tracer{tracer: t}
}

type tracer struct {
	tracer opentracing.Tracer
}

func (t tracer) Start(ctx context.Context, jobName string) (context.Context, cron.TraceFinisher) {
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeyResourceName, Value: jobName},
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeJob},
		tracing.TagMeasured(),
		opentracing.StartTime(time.Now().UTC()),
	}
	span, spanCtx := tracing.StartSpanFromContext(ctx, t.tracer, tracing.OperationJob, startOptions...)
	return spanCtx, &traceFinisher{span: span}
}

type traceFinisher struct {
	span opentracing.Span
}

func (tf traceFinisher) Finish(ctx context.Context, err error) {
	if tf.span == nil {
		return
	}
	tracing.SpanError(tf.span, err)
	tf.span.Finish()
}
