package dbtrace

import (
	"context"
	"database/sql/driver"
	"fmt"
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/tracing"
)

func TestPrepare(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	dbTracer := Tracer(mockTracer)

	dbCfg, err := defaultDB().Config.Reparse()
	assert.Nil(err)

	dbtf := dbTracer.Prepare(context.Background(), dbCfg, "select * from test_table limit 1")
	span := dbtf.(dbTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(tracing.OperationSQLPrepare, mockSpan.OperationName)

	assert.Len(mockSpan.Tags(), 4)
	assert.Equal(tracing.SpanTypeSQL, mockSpan.Tags()[tracing.TagKeySpanType])
	assert.Equal(dbCfg.Database, mockSpan.Tags()[tracing.TagKeyDBName])
	assert.Equal(dbCfg.Username, mockSpan.Tags()[tracing.TagKeyDBUser])
	assert.Equal("select * from test_table limit 1", mockSpan.Tags()[TagKeyQuery])
	assert.True(mockSpan.FinishTime.IsZero())
}

func TestPrepareWithParentSpan(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	dbTracer := Tracer(mockTracer)

	dbCfg, err := defaultDB().Config.Reparse()
	assert.Nil(err)

	parentSpan := mockTracer.StartSpan("test_op")
	ctx := opentracing.ContextWithSpan(context.Background(), parentSpan)

	dbtf := dbTracer.Prepare(ctx, dbCfg, "select * from test_table limit 1")
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

	dbCfg, err := defaultDB().Config.Reparse()
	assert.Nil(err)

	statement := "SELECT 1 FROM test_table WHERE id = $1"
	invocation := defaultDB().Invoke()
	invocation.Label = "test_table_exists"

	dbtf := dbTracer.Query(context.Background(), dbCfg, invocation.Label, statement)
	span := dbtf.(dbTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Equal(tracing.OperationSQLQuery, mockSpan.OperationName)

	assert.Len(mockSpan.Tags(), 5)
	assert.Equal("test_table_exists", mockSpan.Tags()[tracing.TagKeyResourceName])
	assert.Equal(tracing.SpanTypeSQL, mockSpan.Tags()[tracing.TagKeySpanType])
	assert.Equal(dbCfg.Database, mockSpan.Tags()[tracing.TagKeyDBName])
	assert.Equal(dbCfg.Username, mockSpan.Tags()[tracing.TagKeyDBUser])
	assert.Equal(statement, mockSpan.Tags()[TagKeyQuery])
	assert.True(mockSpan.FinishTime.IsZero())
}

func TestQueryWithParentSpan(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	dbTracer := Tracer(mockTracer)

	dbCfg, err := defaultDB().Config.Reparse()
	assert.Nil(err)

	parentSpan := mockTracer.StartSpan("test_op")
	ctx := opentracing.ContextWithSpan(context.Background(), parentSpan)

	statement := "SELECT 1 FROM test_table WHERE id = $1"
	invocation := defaultDB().Invoke()
	invocation.Label = "test_table_exists"

	dbtf := dbTracer.Query(ctx, dbCfg, invocation.Label, statement)
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

	dbCfg, err := defaultDB().Config.Reparse()
	assert.Nil(err)

	dbtf := dbTracer.Query(context.Background(), dbCfg, "ok", "select 'ok1'")
	dbtf.FinishQuery(context.TODO(), nil, nil)

	span := dbtf.(dbTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Nil(mockSpan.Tags()[tracing.TagKeyError])
	assert.False(mockSpan.FinishTime.IsZero())
}

func TestFinishPrepare(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	dbTracer := Tracer(mockTracer)

	dbCfg, err := defaultDB().Config.Reparse()
	assert.Nil(err)

	dbtf := dbTracer.Prepare(context.Background(), dbCfg, "select 'ok1'")
	dbtf.FinishPrepare(context.TODO(), nil)

	span := dbtf.(dbTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Nil(mockSpan.Tags()[tracing.TagKeyError])
	assert.False(mockSpan.FinishTime.IsZero())
}

func TestFinishQueryError(t *testing.T) {
	assert := assert.New(t)
	mockTracer := mocktracer.New()
	dbTracer := Tracer(mockTracer)

	dbCfg, err := defaultDB().Config.Reparse()
	assert.Nil(err)

	ctx := context.Background()
	dbtf := dbTracer.Query(ctx, dbCfg, "ok", "select 'ok1'")
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

	dbCfg, err := defaultDB().Config.Reparse()
	assert.Nil(err)

	ctx := context.Background()
	dbtf := dbTracer.Prepare(ctx, dbCfg, "select 'ok1'")
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

	dbCfg, err := defaultDB().Config.Reparse()
	assert.Nil(err)

	ctx := context.Background()
	dbtf := dbTracer.Query(ctx, dbCfg, "ok", "select 'ok1'")
	dbtf.FinishQuery(ctx, nil, driver.ErrSkip)

	span := dbtf.(dbTraceFinisher).span
	mockSpan := span.(*mocktracer.MockSpan)
	assert.Nil(mockSpan.Tags()[tracing.TagKeyError])
}

func TestFinishQueryNil(t *testing.T) {
	assert := assert.New(t)

	dbtf := dbTraceFinisher{}
	dbtf.FinishQuery(context.TODO(), nil, nil)
	assert.Nil(dbtf.span)
}
