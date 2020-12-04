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
	"github.com/blend/go-sdk/grpcutil"
	"github.com/blend/go-sdk/grpcutil/calculator"
	"github.com/blend/go-sdk/tracing"

	v1 "github.com/blend/go-sdk/grpcutil/calculator/v1"
)

func Test_Tracing_ServerUnary(t *testing.T) {
	assert := assert.New(t)

	mockTracer := mocktracer.New()
	tracer := Tracer(mockTracer)

	// start mocked server with tracing enabled
	socketListener, err := net.Listen("tcp", "127.0.0.1:")
	assert.Nil(err)
	defer socketListener.Close()

	server := grpc.NewServer(grpc.UnaryInterceptor(grpcutil.TracedServerUnary(tracer)))
	v1.RegisterCalculatorServer(server, new(calculator.Server))
	go func() { _ = server.Serve(socketListener) }()

	conn, err := grpc.Dial(socketListener.Addr().String(), grpc.WithInsecure())
	assert.Nil(err)
	res, err := v1.NewCalculatorClient(conn).Add(context.Background(), &v1.Numbers{Values: []float64{1, 2, 3, 4}})
	assert.Nil(err)
	assert.Equal(10, res.Value)

	assert.Len(mockTracer.FinishedSpans(), 1)
	assert.Equal("rpc", mockTracer.FinishedSpans()[0].OperationName)
	assert.Equal("/v1.Calculator/Add", mockTracer.FinishedSpans()[0].Tags()[tracing.TagKeyResourceName])
	assert.Equal("server", mockTracer.FinishedSpans()[0].Tags()[tracing.TagKeyGRPCRole])
	assert.Equal("unary", mockTracer.FinishedSpans()[0].Tags()[tracing.TagKeyGRPCCallingConvention])
}

func Test_Tracing_ServerStream(t *testing.T) {
	assert := assert.New(t)

	mockTracer := mocktracer.New()
	tracer := Tracer(mockTracer)

	// start mocked server with tracing enabled
	socketListener, err := net.Listen("tcp", "127.0.0.1:")
	assert.Nil(err)
	defer socketListener.Close()

	server := grpc.NewServer(grpc.StreamInterceptor(grpcutil.TracedServerStream(tracer)))
	v1.RegisterCalculatorServer(server, new(calculator.Server))
	go func() { _ = server.Serve(socketListener) }()

	conn, err := grpc.Dial(socketListener.Addr().String(), grpc.WithInsecure())
	assert.Nil(err)
	stream, err := v1.NewCalculatorClient(conn).AddStream(context.Background())
	assert.Nil(err)

	assert.Nil(stream.Send(&v1.Number{Value: 1}))
	assert.Nil(stream.Send(&v1.Number{Value: 2}))
	assert.Nil(stream.Send(&v1.Number{Value: 3}))
	assert.Nil(stream.Send(&v1.Number{Value: 4}))

	res, err := stream.CloseAndRecv()
	assert.Nil(err)
	assert.Equal(10, res.Value)

	assert.Len(mockTracer.FinishedSpans(), 1)
	assert.Equal("rpc", mockTracer.FinishedSpans()[0].OperationName)
	assert.Equal("/v1.Calculator/AddStream", mockTracer.FinishedSpans()[0].Tags()[tracing.TagKeyResourceName])
	assert.Equal("server", mockTracer.FinishedSpans()[0].Tags()[tracing.TagKeyGRPCRole])
	assert.Equal("stream", mockTracer.FinishedSpans()[0].Tags()[tracing.TagKeyGRPCCallingConvention])
}

func Test_Tracing_ClientServerUnary(t *testing.T) {
	assert := assert.New(t)

	mockTracer := mocktracer.New()
	tracer := Tracer(mockTracer)

	// start mocked server with tracing enabled
	socketListener, err := net.Listen("tcp", "127.0.0.1:")
	assert.Nil(err)
	defer socketListener.Close()

	server := grpc.NewServer(grpc.UnaryInterceptor(grpcutil.TracedServerUnary(tracer)))
	v1.RegisterCalculatorServer(server, new(calculator.Server))
	go func() { _ = server.Serve(socketListener) }()

	conn, err := grpc.Dial(socketListener.Addr().String(), grpc.WithInsecure(), grpc.WithUnaryInterceptor(grpcutil.TracedClientUnary(tracer)))
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
	assert.Equal("unary", mockTracer.FinishedSpans()[0].Tags()[tracing.TagKeyGRPCCallingConvention])

	// client
	assert.Zero(mockTracer.FinishedSpans()[1].ParentID)
	assert.Equal("rpc", mockTracer.FinishedSpans()[1].OperationName)
	assert.Equal("/v1.Calculator/Add", mockTracer.FinishedSpans()[1].Tags()[tracing.TagKeyResourceName])
	assert.Equal("client", mockTracer.FinishedSpans()[1].Tags()[tracing.TagKeyGRPCRole])
	assert.Equal("unary", mockTracer.FinishedSpans()[1].Tags()[tracing.TagKeyGRPCCallingConvention])
}

func Test_Tracing_ClientServerStream(t *testing.T) {
	assert := assert.New(t)

	mockTracer := mocktracer.New()
	tracer := Tracer(mockTracer)

	// start mocked server with tracing enabled
	socketListener, err := net.Listen("tcp", "127.0.0.1:")
	assert.Nil(err)
	defer socketListener.Close()

	server := grpc.NewServer(grpc.StreamInterceptor(grpcutil.TracedServerStream(tracer)))
	v1.RegisterCalculatorServer(server, new(calculator.Server))
	go func() { _ = server.Serve(socketListener) }()

	conn, err := grpc.Dial(socketListener.Addr().String(), grpc.WithInsecure(), grpc.WithStreamInterceptor(grpcutil.TracedClientStream(tracer)))
	assert.Nil(err)
	stream, err := v1.NewCalculatorClient(conn).AddStream(context.Background())
	assert.Nil(err)

	assert.Nil(stream.Send(&v1.Number{Value: 1}))
	assert.Nil(stream.Send(&v1.Number{Value: 2}))
	assert.Nil(stream.Send(&v1.Number{Value: 3}))
	assert.Nil(stream.Send(&v1.Number{Value: 4}))

	res, err := stream.CloseAndRecv()
	assert.Nil(err)
	assert.Equal(10, res.Value)

	assert.Len(mockTracer.FinishedSpans(), 2)
	assert.Equal("rpc", mockTracer.FinishedSpans()[0].OperationName)
	assert.Equal("/v1.Calculator/AddStream", mockTracer.FinishedSpans()[0].Tags()[tracing.TagKeyResourceName])
	assert.Equal("client", mockTracer.FinishedSpans()[0].Tags()[tracing.TagKeyGRPCRole])
	assert.Equal("stream", mockTracer.FinishedSpans()[0].Tags()[tracing.TagKeyGRPCCallingConvention])

	assert.NotZero(mockTracer.FinishedSpans()[1].ParentID)
	assert.Equal("rpc", mockTracer.FinishedSpans()[1].OperationName)
	assert.Equal("/v1.Calculator/AddStream", mockTracer.FinishedSpans()[1].Tags()[tracing.TagKeyResourceName])
	assert.Equal("server", mockTracer.FinishedSpans()[1].Tags()[tracing.TagKeyGRPCRole])
	assert.Equal("stream", mockTracer.FinishedSpans()[1].Tags()[tracing.TagKeyGRPCCallingConvention])
}

func Test_Tracing_ParentClientServerUnary(t *testing.T) {
	assert := assert.New(t)

	mockTracer := mocktracer.New()
	tracer := Tracer(mockTracer)

	// start mocked server with tracing enabled
	socketListener, err := net.Listen("tcp", "127.0.0.1:")
	assert.Nil(err)
	defer socketListener.Close()

	server := grpc.NewServer(grpc.UnaryInterceptor(grpcutil.TracedServerUnary(tracer)))
	v1.RegisterCalculatorServer(server, new(calculator.Server))
	go func() { _ = server.Serve(socketListener) }()

	outerSpan, ctx := tracing.StartSpanFromContext(context.Background(), mockTracer, tracing.OperationHTTPRequest,
		opentracing.Tag{Key: tracing.TagKeyResourceName, Value: "/foo"},
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeWeb},
		opentracing.StartTime(time.Now().UTC()),
	)

	conn, err := grpc.Dial(socketListener.Addr().String(), grpc.WithInsecure(), grpc.WithUnaryInterceptor(grpcutil.TracedClientUnary(tracer)))
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
	assert.Equal("unary", mockTracer.FinishedSpans()[0].Tags()[tracing.TagKeyGRPCCallingConvention])

	// client
	assert.NotZero(mockTracer.FinishedSpans()[1].ParentID)
	assert.Equal("rpc", mockTracer.FinishedSpans()[1].OperationName)
	assert.Equal("/v1.Calculator/Add", mockTracer.FinishedSpans()[1].Tags()[tracing.TagKeyResourceName])
	assert.Equal("client", mockTracer.FinishedSpans()[1].Tags()[tracing.TagKeyGRPCRole])
	assert.Equal("unary", mockTracer.FinishedSpans()[1].Tags()[tracing.TagKeyGRPCCallingConvention])
}

func Test_Tracing_ParentClientServerStream(t *testing.T) {
	assert := assert.New(t)

	mockTracer := mocktracer.New()
	tracer := Tracer(mockTracer)

	// start mocked server with tracing enabled
	socketListener, err := net.Listen("tcp", "127.0.0.1:")
	assert.Nil(err)
	defer socketListener.Close()

	server := grpc.NewServer(grpc.StreamInterceptor(grpcutil.TracedServerStream(tracer)))
	v1.RegisterCalculatorServer(server, new(calculator.Server))
	go func() { _ = server.Serve(socketListener) }()

	conn, err := grpc.Dial(socketListener.Addr().String(), grpc.WithInsecure(), grpc.WithStreamInterceptor(grpcutil.TracedClientStream(tracer)))
	assert.Nil(err)

	outerSpan, ctx := tracing.StartSpanFromContext(context.Background(), mockTracer, tracing.OperationHTTPRequest,
		opentracing.Tag{Key: tracing.TagKeyResourceName, Value: "/foo"},
		opentracing.Tag{Key: tracing.TagKeySpanType, Value: tracing.SpanTypeWeb},
		opentracing.StartTime(time.Now().UTC()),
	)

	stream, err := v1.NewCalculatorClient(conn).AddStream(ctx)
	assert.Nil(err)

	assert.Nil(stream.Send(&v1.Number{Value: 1}))
	assert.Nil(stream.Send(&v1.Number{Value: 2}))
	assert.Nil(stream.Send(&v1.Number{Value: 3}))
	assert.Nil(stream.Send(&v1.Number{Value: 4}))

	res, err := stream.CloseAndRecv()
	assert.Nil(err)
	assert.Equal(10, res.Value)

	// finish the outer span ...
	outerSpan.Finish()

	assert.Len(mockTracer.FinishedSpans(), 3)

	assert.NotZero(mockTracer.FinishedSpans()[0].ParentID)
	assert.Equal("rpc", mockTracer.FinishedSpans()[0].OperationName)
	assert.Equal("/v1.Calculator/AddStream", mockTracer.FinishedSpans()[0].Tags()[tracing.TagKeyResourceName])
	assert.Equal("client", mockTracer.FinishedSpans()[0].Tags()[tracing.TagKeyGRPCRole])
	assert.Equal("stream", mockTracer.FinishedSpans()[0].Tags()[tracing.TagKeyGRPCCallingConvention])

	assert.NotZero(mockTracer.FinishedSpans()[1].ParentID)
	assert.Equal("rpc", mockTracer.FinishedSpans()[1].OperationName)
	assert.Equal("/v1.Calculator/AddStream", mockTracer.FinishedSpans()[1].Tags()[tracing.TagKeyResourceName])
	assert.Equal("server", mockTracer.FinishedSpans()[1].Tags()[tracing.TagKeyGRPCRole])
	assert.Equal("stream", mockTracer.FinishedSpans()[1].Tags()[tracing.TagKeyGRPCCallingConvention])
}
