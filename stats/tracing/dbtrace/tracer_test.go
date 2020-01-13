package dbtrace

import (
	"context"
	"database/sql/driver"
	"fmt"
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/stats/tracing"
)

func TestPrepare(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	dbTracer := Tracer(mockTracer)

	dbtf := dbTracer.Prepare(context.Background(), defaultDB().Config, "select * from test_table limit 1")
	span := dbtf.(dbTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(tracing.OperationSQLPrepare, mockSpan.OperationName)

	assert.Len(mockSpan.Tags(), 4)
	assert.Equal(tracing.SpanTypeSQL, mockSpan.Tags()[tracing.TagKeySpanType])
	assert.Equal("postgres", mockSpan.Tags()[tracing.TagKeyDBName])
	assert.Equal("", mockSpan.Tags()[tracing.TagKeyDBUser])
	assert.Equal("select * from test_table limit 1", mockSpan.Tags()[TagKeyQuery])
	assert.True(mockSpan.FinishTime.IsZero())
}

func TestPrepareWithParentSpan(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	dbTracer := Tracer(mockTracer)

	parentSpan := mockTracer.StartSpan("test_op")
	ctx := opentracing.ContextWithSpan(context.Background(), parentSpan)

	dbtf := dbTracer.Prepare(ctx, defaultDB().Config, "select * from test_table limit 1")
	span := dbtf.(dbTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(tracing.OperationSQLPrepare, mockSpan.OperationName)

	mockParentSpan := parentSpan.(*mocktracer.MockSpan)
	assert.Equal(mockSpan.ParentID, mockParentSpan.SpanContext.SpanID)
}

func TestQuery(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	dbTracer := Tracer(mockTracer)

	statement := "SELECT 1 FROM test_table WHERE id = $1"
	invocation := defaultDB().Invoke()
	invocation.Label = "test_table_exists"

	dbtf := dbTracer.Query(context.Background(), defaultDB().Config, invocation.Label, statement)
	span := dbtf.(dbTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(tracing.OperationSQLQuery, mockSpan.OperationName)

	assert.Len(mockSpan.Tags(), 5)
	assert.Equal("test_table_exists", mockSpan.Tags()[tracing.TagKeyResourceName])
	assert.Equal(tracing.SpanTypeSQL, mockSpan.Tags()[tracing.TagKeySpanType])
	assert.Equal("postgres", mockSpan.Tags()[tracing.TagKeyDBName])
	assert.Equal("", mockSpan.Tags()[tracing.TagKeyDBUser])
	assert.Equal(statement, mockSpan.Tags()[TagKeyQuery])
	assert.True(mockSpan.FinishTime.IsZero())
}

func TestQueryWithParentSpan(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	dbTracer := Tracer(mockTracer)

	parentSpan := mockTracer.StartSpan("test_op")
	ctx := opentracing.ContextWithSpan(context.Background(), parentSpan)

	statement := "SELECT 1 FROM test_table WHERE id = $1"
	invocation := defaultDB().Invoke()
	invocation.Label = "test_table_exists"

	dbtf := dbTracer.Query(ctx, defaultDB().Config, invocation.Label, statement)
	span := dbtf.(dbTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(tracing.OperationSQLQuery, mockSpan.OperationName)

	mockParentSpan := parentSpan.(*mocktracer.MockSpan)
	assert.Equal(mockSpan.ParentID, mockParentSpan.SpanContext.SpanID)
}

func TestFinishQuery(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	dbTracer := Tracer(mockTracer)

	dbtf := dbTracer.Query(context.Background(), defaultDB().Config, "ok", "select 'ok1'")
	dbtf.FinishQuery(nil, nil, nil)

	span := dbtf.(dbTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Nil(mockSpan.Tags()[tracing.TagKeyError])
	assert.False(mockSpan.FinishTime.IsZero())
}

func TestFinishPrepare(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	dbTracer := Tracer(mockTracer)

	dbtf := dbTracer.Prepare(context.Background(), defaultDB().Config, "select 'ok1'")
	dbtf.FinishPrepare(nil, nil)

	span := dbtf.(dbTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Nil(mockSpan.Tags()[tracing.TagKeyError])
	assert.False(mockSpan.FinishTime.IsZero())
}

func TestFinishQueryError(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	dbTracer := Tracer(mockTracer)

	ctx := context.Background()
	dbtf := dbTracer.Query(ctx, defaultDB().Config, "ok", "select 'ok1'")
	dbtf.FinishQuery(ctx, nil, fmt.Errorf("error"))

	span := dbtf.(dbTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal("error", mockSpan.Tags()[tracing.TagKeyError])
	assert.False(mockSpan.FinishTime.IsZero())
}

func TestFinishPrepareError(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	dbTracer := Tracer(mockTracer)

	ctx := context.Background()
	dbtf := dbTracer.Prepare(ctx, defaultDB().Config, "select 'ok1'")
	dbtf.FinishPrepare(ctx, fmt.Errorf("error"))

	span := dbtf.(dbTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal("error", mockSpan.Tags()[tracing.TagKeyError])
	assert.False(mockSpan.FinishTime.IsZero())
}

func TestFinishQueryErrorSkip(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	dbTracer := Tracer(mockTracer)

	ctx := context.Background()
	dbtf := dbTracer.Query(ctx, defaultDB().Config, "ok", "select 'ok1'")
	dbtf.FinishQuery(ctx, nil, driver.ErrSkip)

	span := dbtf.(dbTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Nil(mockSpan.Tags()[tracing.TagKeyError])
}

func TestFinishQueryNil(t *testing.T) {
	assert := assert.New(t)

	dbtf := dbTraceFinisher{}
	dbtf.FinishQuery(nil, nil, nil)
	assert.Nil(dbtf.span)
}
