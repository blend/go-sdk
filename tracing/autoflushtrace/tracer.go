package autoflushtrace

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/blend/go-sdk/autoflush"
	"github.com/blend/go-sdk/tracing"
)

var (
	_ autoflush.Tracer        = (*autoflushTracer)(nil)
	_ autoflush.TraceFinisher = (*autoflushTraceFinisher)(nil)
)

// Tracer returns a new tracer.
func Tracer(tracer opentracing.Tracer) autoflush.Tracer {
	return &autoflushTracer{tracer: tracer}
}

type autoflushTracer struct {
	tracer opentracing.Tracer
}

func (aft autoflushTracer) StartAdd(ctx context.Context) autoflush.TraceFinisher {
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: "queue"},
		tracing.TagMeasured(),
		opentracing.StartTime(time.Now().UTC()),
	}
	span, _ := tracing.StartSpanFromContext(ctx, aft.tracer, "autoflush.add", startOptions...)
	return autoflushTraceFinisher{span: span}
}

func (aft autoflushTracer) StartAddMany(ctx context.Context) autoflush.TraceFinisher {
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: "queue"},
		tracing.TagMeasured(),
		opentracing.StartTime(time.Now().UTC()),
	}
	span, _ := tracing.StartSpanFromContext(ctx, aft.tracer, "autoflush.add_many", startOptions...)
	return autoflushTraceFinisher{span: span}
}

func (aft autoflushTracer) StartQueueFlush(ctx context.Context) autoflush.TraceFinisher {
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: "queue"},
		tracing.TagMeasured(),
		opentracing.StartTime(time.Now().UTC()),
	}
	span, _ := tracing.StartSpanFromContext(ctx, aft.tracer, "autoflush.queue_flush", startOptions...)
	return autoflushTraceFinisher{span: span}
}

func (aft autoflushTracer) StartFlush(ctx context.Context) (context.Context, autoflush.TraceFinisher) {
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: "queue"},
		tracing.TagMeasured(),
		opentracing.StartTime(time.Now().UTC()),
	}
	span, ctx := tracing.StartSpanFromContext(ctx, aft.tracer, "autoflush.flush", startOptions...)
	return ctx, autoflushTraceFinisher{span: span}
}

type autoflushTraceFinisher struct {
	span opentracing.Span
}

func (tf autoflushTraceFinisher) Finish(err error) {
	if err != nil {
		tracing.SpanError(tf.span, err)
	}
	tf.span.Finish()
}
