package crontrace

import (
	"context"
	"time"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/stats/tracing"
	opentracing "github.com/opentracing/opentracing-go"
)

// Tracer returns a opentracing cron tracer.
func Tracer(t opentracing.Tracer) cron.Tracer {
	return &tracer{tracer: t}
}

type tracer struct {
	tracer opentracing.Tracer
}

func (t tracer) Start(ctx context.Context) (context.Context, cron.TraceFinisher) {
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeyResourceName, Value: cron.GetJobInvocation(ctx).JobName},
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeJob},
		opentracing.StartTime(time.Now().UTC()),
	}
	span, spanCtx := tracing.StartSpanFromContext(ctx, t.tracer, tracing.OperationJob, startOptions...)
	return spanCtx, &traceFinisher{span: span}
}

type traceFinisher struct {
	span opentracing.Span
}

func (tf traceFinisher) Finish(ctx context.Context) {
	if tf.span == nil {
		return
	}
	tracing.SpanError(tf.span, cron.GetJobInvocation(ctx).Err)
	tf.span.Finish()
}
