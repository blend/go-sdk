package crontrace

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/stats/tracing"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
)

func TestStart(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	cronTracer := Tracer(mockTracer)

	ji := cron.NewJobInvocation("test_job")
	ctx := cron.WithJobInvocation(context.Background(), ji)
	ctx, _ = cronTracer.Start(ctx)

	span := opentracing.SpanFromContext(ctx)
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(tracing.OperationJob, mockSpan.OperationName)

	assert.Len(mockSpan.Tags(), 2)
	assert.Equal("test_job", mockSpan.Tags()[tracing.TagKeyResourceName])
	assert.Equal(tracing.SpanTypeJob, mockSpan.Tags()[tracing.TagKeySpanType])
	assert.True(mockSpan.FinishTime.IsZero())
}

func TestFinish(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	cronTracer := Tracer(mockTracer)
	testStartTime := time.Now()

	ji := cron.NewJobInvocation("test_job")
	ctx := cron.WithJobInvocation(context.Background(), ji)
	ctx, tf := cronTracer.Start(ctx)

	tf.Finish(ctx)
	span := opentracing.SpanFromContext(ctx)
	mockSpan := span.(*mocktracer.MockSpan)
	assert.True(testStartTime.Before(mockSpan.FinishTime))
	assert.Equal(nil, mockSpan.Tags()[tracing.TagKeyError])
	assert.False(mockSpan.FinishTime.IsZero())
}

func TestFinishError(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	cronTracer := Tracer(mockTracer)
	testStartTime := time.Now()

	ji := cron.NewJobInvocation("test_job")
	ji.Err = fmt.Errorf("error")

	ctx := cron.WithJobInvocation(context.Background(), ji)
	// Start Span from Background Context
	ctx, tf := cronTracer.Start(ctx)

	tf.Finish(ctx)
	span := opentracing.SpanFromContext(ctx)
	mockSpan := span.(*mocktracer.MockSpan)
	assert.True(testStartTime.Before(mockSpan.FinishTime))
	assert.Equal("error", mockSpan.Tags()[tracing.TagKeyError])
	assert.False(mockSpan.FinishTime.IsZero())
}

func TestFinishNilSpan(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	traceFinisher{}.Finish(ctx)
	assert.Nil(opentracing.SpanFromContext(ctx))
}
