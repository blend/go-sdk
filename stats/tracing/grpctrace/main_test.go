package grpctrace

import (
	"context"
	"net"
	"testing"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"google.golang.org/grpc"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/grpcutil/calculator"
	v1 "github.com/blend/go-sdk/grpcutil/calculator/v1"
	"github.com/blend/go-sdk/stats/tracing"
)

func Test_Tracing_ServerUnary(t *testing.T) {
	assert := assert.New(t)

	mockTracer := mocktracer.New()

	// start mocked server with tracing enabled
	socketListener, err := net.Listen("tcp", "127.0.0.1:")
	assert.Nil(err)
	defer socketListener.Close()

	server := grpc.NewServer(grpc.UnaryInterceptor(TracedServerUnary(mockTracer)))
	v1.RegisterCalculatorServer(server, new(calculator.Server))
	go server.Serve(socketListener)

	conn, err := grpc.Dial(socketListener.Addr().String(), grpc.WithInsecure())
	assert.Nil(err)
	res, err := v1.NewCalculatorClient(conn).Add(context.Background(), &v1.Numbers{Values: []float64{1, 2, 3, 4}})
	assert.Nil(err)
	assert.Equal(10, res.Value)

	assert.Len(mockTracer.FinishedSpans(), 1)
	assert.Equal("rpc", mockTracer.FinishedSpans()[0].OperationName)
	assert.Equal("/v1.Calculator/Add", mockTracer.FinishedSpans()[0].Tags()["resource.name"])
}

func Test_Tracing_ClientServerUnary(t *testing.T) {
	assert := assert.New(t)

	mockTracer := mocktracer.New()

	// start mocked server with tracing enabled
	socketListener, err := net.Listen("tcp", "127.0.0.1:")
	assert.Nil(err)
	defer socketListener.Close()

	server := grpc.NewServer(grpc.UnaryInterceptor(TracedServerUnary(mockTracer)))
	v1.RegisterCalculatorServer(server, new(calculator.Server))
	go server.Serve(socketListener)

	conn, err := grpc.Dial(socketListener.Addr().String(), grpc.WithInsecure(), grpc.WithUnaryInterceptor(TracedClientUnary(mockTracer)))
	assert.Nil(err)
	res, err := v1.NewCalculatorClient(conn).Add(context.Background(), &v1.Numbers{Values: []float64{1, 2, 3, 4}})
	assert.Nil(err)
	assert.Equal(10, res.Value)

	assert.Len(mockTracer.FinishedSpans(), 2)

	// server
	assert.NotZero(mockTracer.FinishedSpans()[0].ParentID)
	assert.Equal("rpc", mockTracer.FinishedSpans()[0].OperationName)
	assert.Equal("/v1.Calculator/Add", mockTracer.FinishedSpans()[0].Tags()[tracing.TagKeyResourceName])
	assert.Equal("server", mockTracer.FinishedSpans()[0].Tags()[tracing.TagKeyGRPCRole])

	// client
	assert.Zero(mockTracer.FinishedSpans()[1].ParentID)
	assert.Equal("rpc", mockTracer.FinishedSpans()[1].OperationName)
	assert.Equal("/v1.Calculator/Add", mockTracer.FinishedSpans()[1].Tags()[tracing.TagKeyResourceName])
	assert.Equal("client", mockTracer.FinishedSpans()[1].Tags()[tracing.TagKeyGRPCRole])
}

func Test_Tracing_ParentClientServerUnary(t *testing.T) {
	assert := assert.New(t)

	mockTracer := mocktracer.New()

	// start mocked server with tracing enabled
	socketListener, err := net.Listen("tcp", "127.0.0.1:")
	assert.Nil(err)
	defer socketListener.Close()

	server := grpc.NewServer(grpc.UnaryInterceptor(TracedServerUnary(mockTracer)))
	v1.RegisterCalculatorServer(server, new(calculator.Server))
	go server.Serve(socketListener)

	outerSpan, ctx := tracing.StartSpanFromContext(context.Background(), mockTracer, tracing.OperationHTTPRequest,
		opentracing.Tag{Key: tracing.TagKeyResourceName, Value: "/foo"},
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeWeb},
		opentracing.StartTime(time.Now().UTC()),
	)

	conn, err := grpc.Dial(socketListener.Addr().String(), grpc.WithInsecure(), grpc.WithUnaryInterceptor(TracedClientUnary(mockTracer)))
	assert.Nil(err)
	res, err := v1.NewCalculatorClient(conn).Add(ctx, &v1.Numbers{Values: []float64{1, 2, 3, 4}})
	assert.Nil(err)
	assert.Equal(10, res.Value)

	// finish the outer span ...
	outerSpan.Finish()

	assert.Len(mockTracer.FinishedSpans(), 3)

	// server
	assert.NotZero(mockTracer.FinishedSpans()[0].ParentID)
	assert.Equal("rpc", mockTracer.FinishedSpans()[0].OperationName)
	assert.Equal("/v1.Calculator/Add", mockTracer.FinishedSpans()[0].Tags()[tracing.TagKeyResourceName])
	assert.Equal("server", mockTracer.FinishedSpans()[0].Tags()[tracing.TagKeyGRPCRole])

	// client
	assert.NotZero(mockTracer.FinishedSpans()[1].ParentID)
	assert.Equal("rpc", mockTracer.FinishedSpans()[0].OperationName)
	assert.Equal("/v1.Calculator/Add", mockTracer.FinishedSpans()[1].Tags()[tracing.TagKeyResourceName])
	assert.Equal("client", mockTracer.FinishedSpans()[1].Tags()[tracing.TagKeyGRPCRole])
}
