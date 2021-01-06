package grpcutil

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/grpc"

	"github.com/blend/go-sdk/assert"
)

func TestRecoverUnary(t *testing.T) {
	assert := assert.New(t)

	interceptor := RecoverServerUnary(WithServerRecoveryHandler(func(p interface{}) error {
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

	interceptor := RecoverServerStream(WithServerRecoveryHandler(func(p interface{}) error {
		return fmt.Errorf("panic: %v", p)
	}))

	err := interceptor(nil, nil, nil, func(srv interface{}, stream grpc.ServerStream) error {
		panic("errored in handler")
	})
	assert.NotNil(err)
	assert.Equal("panic: errored in handler", err.Error())
}
