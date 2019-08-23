package dbtrace

import (
	"context"
	"database/sql/driver"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/db"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/stats/tracing"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
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

	dbtf := dbTracer.Query(context.Background(), defaultDB().Config, invocation, statement)
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

	dbtf := dbTracer.Query(ctx, defaultDB().Config, invocation, statement)
	span := dbtf.(dbTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(tracing.OperationSQLQuery, mockSpan.OperationName)

	mockParentSpan := parentSpan.(*mocktracer.MockSpan)
	assert.Equal(mockSpan.ParentID, mockParentSpan.SpanContext.SpanID)
}

func TestFinish(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	dbTracer := Tracer(mockTracer)

	dbtf := dbTracer.Query(context.Background(), defaultDB().Config, &db.Invocation{}, "select 'ok1'")
	dbtf.Finish(nil)

	span := dbtf.(dbTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Nil(mockSpan.Tags()[tracing.TagKeyError])
	assert.False(mockSpan.FinishTime.IsZero())
}

func TestFinishError(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	dbTracer := Tracer(mockTracer)

	dbtf := dbTracer.Query(context.Background(), defaultDB().Config, &db.Invocation{}, "select 'ok1'")
	dbtf.Finish(fmt.Errorf("error"))

	span := dbtf.(dbTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal("error", mockSpan.Tags()[tracing.TagKeyError])
	assert.False(mockSpan.FinishTime.IsZero())
}

func TestFinishErrorSkip(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	dbTracer := Tracer(mockTracer)

	dbtf := dbTracer.Query(context.Background(), defaultDB().Config, &db.Invocation{}, "select 'ok1'")
	dbtf.Finish(driver.ErrSkip)

	span := dbtf.(dbTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Nil(mockSpan.Tags()[tracing.TagKeyError])
}

func TestFinishNil(t *testing.T) {
	assert := assert.New(t)

	dbtf := dbTraceFinisher{}
	dbtf.Finish(nil)
	assert.Nil(dbtf.span)
}
