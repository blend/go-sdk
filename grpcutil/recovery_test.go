package grpcutil

import (
	"context"
	"fmt"
	"testing"
	"github.com/blend/go-sdk/assert"
	"google.golang.org/grpc"
)

func TestRecoverUnary(t *testing.T) {
	assert := assert.New(t)

	interceptor := RecoverUnary(WithRecoveryHandler(func(p interface{}) error {
		return fmt.Errorf("panic: %v", p)
	}))

	_, err := interceptor(context.TODO(), nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("errored in handler")
	})
	assert.NotNil(err)
	assert.Equal("panic: errored in handler", err.Error())
}

func TestRecoverStream(t *testing.T) {
	assert := assert.New(t)

	interceptor := RecoverStream(WithRecoveryHandler(func(p interface{}) error {
		return fmt.Errorf("panic: %v", p)
	}))

	err := interceptor(nil, nil, nil, func(srv interface{}, stream grpc.ServerStream) error {
		panic("errored in handler")
	})
	assert.NotNil(err)
	assert.Equal("panic: errored in handler", err.Error())
}
