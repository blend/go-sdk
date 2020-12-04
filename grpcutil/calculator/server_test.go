package calculator

import (
	"context"
	"io"
	"testing"

	"google.golang.org/grpc/metadata"

	"github.com/blend/go-sdk/assert"
	v1 "github.com/blend/go-sdk/grpcutil/calculator/v1"
)

func numberStream(values ...float64) *numberStreamServer {
	return &numberStreamServer{
		Input: values,
	}
}

var (
	_ v1.Calculator_AddStreamServer      = (*numberStreamServer)(nil)
	_ v1.Calculator_SubtractStreamServer = (*numberStreamServer)(nil)
)

type mockServerStream struct{}

func (mss mockServerStream) SetHeader(metadata.MD) error  { return nil }
func (mss mockServerStream) SetTrailer(metadata.MD)       {}
func (mss mockServerStream) SendHeader(metadata.MD) error { return nil }
func (mss mockServerStream) Context() context.Context     { return context.Background() }
func (mss mockServerStream) SendMsg(m interface{}) error  { return nil }
func (mss mockServerStream) RecvMsg(m interface{}) error  { return nil }

type numberStreamServer struct {
	mockServerStream
	Output     float64
	Input      []float64
	InputIndex int
}

func (nss *numberStreamServer) Recv() (*v1.Number, error) {
	defer func() {
		nss.InputIndex++
	}()

	if nss.InputIndex < len(nss.Input) {
		return &v1.Number{
			Value: nss.Input[nss.InputIndex],
		}, nil
	}
	return nil, io.EOF
}

func (nss *numberStreamServer) SendAndClose(value *v1.Number) error {
	nss.Output = value.Value
	return nil
}

func Test_CalculatorServer_Add(t *testing.T) {
	assert := assert.New(t)

	res, err := new(Server).Add(context.TODO(), &v1.Numbers{
		Values: []float64{1, 2, 3, 4},
	})
	assert.Nil(err)
	assert.Equal(10, res.Value)
}

func Test_CalculatorServer_AddStream(t *testing.T) {
	assert := assert.New(t)

	stream := numberStream(1, 2, 3, 4)
	err := new(Server).AddStream(stream)
	assert.Nil(err)
	assert.Equal(10, stream.Output)
}

func Test_CalculatorServer_Subtract(t *testing.T) {
	assert := assert.New(t)

	res, err := new(Server).Subtract(context.TODO(), &v1.Numbers{
		Values: []float64{1, 2, 3, 4},
	})
	assert.Nil(err)
	assert.Equal(-8, res.Value)
}

func Test_CalculatorServer_SubtractStream(t *testing.T) {
	assert := assert.New(t)

	stream := numberStream(1, 2, 3, 4)
	err := new(Server).SubtractStream(stream)
	assert.Nil(err)
	assert.Equal(-8, stream.Output)
}

func Test_CalculatorServer_Multiply(t *testing.T) {
	assert := assert.New(t)

	res, err := new(Server).Multiply(context.TODO(), &v1.Numbers{
		Values: []float64{1, 2, 3, 4},
	})
	assert.Nil(err)
	assert.Equal(24, res.Value)
}

func Test_CalculatorServer_MultiplyStream(t *testing.T) {
	assert := assert.New(t)

	stream := numberStream(1, 2, 3, 4)
	err := new(Server).MultiplyStream(stream)
	assert.Nil(err)
	assert.Equal(24, stream.Output)
}

func Test_CalculatorServer_Divide(t *testing.T) {
	assert := assert.New(t)

	res, err := new(Server).Divide(context.TODO(), &v1.Numbers{
		Values: []float64{1, 2, 3, 4},
	})
	assert.Nil(err)
	assert.InDelta(0.04167, res.Value, 0.0001)
}

func Test_CalculatorServer_DivideStream(t *testing.T) {
	assert := assert.New(t)

	stream := numberStream(1, 2, 3, 4)
	err := new(Server).DivideStream(stream)
	assert.Nil(err)
	assert.InDelta(0.04167, stream.Output, 0.0001)
}
