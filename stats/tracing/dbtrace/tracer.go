package dbtrace

import (
	"context"
	"database/sql/driver"
	"time"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/stats/tracing"
	opentracing "github.com/opentracing/opentracing-go"
)

var (
	_ db.Tracer = (*dbTracer)(nil)
)

// Tracer returns a db tracer.
func Tracer(tracer opentracing.Tracer) db.Tracer {
	return &dbTracer{tracer: tracer}
}

type dbTracer struct {
	tracer opentracing.Tracer
}

func (dbt dbTracer) Prepare(ctx context.Context, cfg db.Config, statement string) db.TraceFinisher {
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeSQL},
		opentracing.Tag{Key: tracing.TagKeyDBName, Value: cfg.DatabaseOrDefault()},
		opentracing.Tag{Key: tracing.TagKeyDBUser, Value: cfg.Username},
		opentracing.Tag{Key: TagKeyQuery, Value: statement},
		opentracing.StartTime(time.Now().UTC()),
	}
	span, _ := tracing.StartSpanFromContext(ctx, dbt.tracer, tracing.OperationSQLPrepare, startOptions...)
	return dbTraceFinisher{span: span}
}

func (dbt dbTracer) Query(ctx context.Context, cfg db.Config, inv *db.Invocation, statement string) db.TraceFinisher {
	startOptions := []opentracing.StartSpanOption{
		opentracing.Tag{Key: tracing.TagKeyResourceName, Value: inv.Label},
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeSQL},
		opentracing.Tag{Key: tracing.TagKeyDBName, Value: cfg.DatabaseOrDefault()},
		opentracing.Tag{Key: tracing.TagKeyDBUser, Value: cfg.Username},
		opentracing.Tag{Key: TagKeyQuery, Value: statement},
		opentracing.StartTime(inv.StartTime),
	}
	span, _ := tracing.StartSpanFromContext(ctx, dbt.tracer, tracing.OperationSQLQuery, startOptions...)
	return dbTraceFinisher{span: span}
}

type dbTraceFinisher struct {
	span opentracing.Span
}

func (dbtf dbTraceFinisher) Finish(err error) {
	if dbtf.span == nil {
		return
	}
	if err == driver.ErrSkip {
		return
	}

	tracing.SpanError(dbtf.span, err)
	dbtf.span.Finish()
}
