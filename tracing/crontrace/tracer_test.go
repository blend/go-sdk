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
	"fmt"
	"testing"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/tracing"
)

func TestStart(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	cronTracer := Tracer(mockTracer)

	ctx := context.Background()
	ctx, _ = cronTracer.Start(ctx, "test_job")

	span := opentracing.SpanFromContext(ctx)
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(tracing.OperationJob, mockSpan.OperationName)

	assert.Len(mockSpan.Tags(), 3)
	assert.Equal("test_job", mockSpan.Tags()[tracing.TagKeyResourceName])
	assert.Equal(tracing.SpanTypeJob, mockSpan.Tags()[tracing.TagKeySpanType])
	assert.True(mockSpan.FinishTime.IsZero())
}

func TestFinish(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	cronTracer := Tracer(mockTracer)
	testStartTime := time.Now()

	ctx := context.Background()
	ctx, tf := cronTracer.Start(ctx, "tracer-test")

	tf.Finish(ctx, nil)
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

	ctx := context.Background()
	// Start Span from Background Context
	ctx, tf := cronTracer.Start(ctx, "tracer-test")

	tf.Finish(ctx, fmt.Errorf("error"))
	span := opentracing.SpanFromContext(ctx)
	mockSpan := span.(*mocktracer.MockSpan)
	assert.True(testStartTime.Before(mockSpan.FinishTime))
	assert.Equal("error", mockSpan.Tags()[tracing.TagKeyError])
	assert.False(mockSpan.FinishTime.IsZero())
}

func TestFinishNilSpan(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	traceFinisher{}.Finish(ctx, nil)
	assert.Nil(opentracing.SpanFromContext(ctx))
}
