package grpcutil

import (
	"context"
	"testing"

	"google.golang.org/grpc"

	"github.com/blend/go-sdk/assert"
)

func TestStreamChain(t *testing.T) {
	assert := assert.New(t)

	var calls []string
	combined := StreamChain(
		func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			calls = append(calls, "first")
			return handler(srv, ss)
		},
		func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			calls = append(calls, "second")
			return handler(srv, ss)
		},
		func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			calls = append(calls, "third")
			return handler(srv, ss)
		},
	)

	err := combined(context.Background(), nil, nil, func(isrv interface{}, istream grpc.ServerStream) error {
		calls = append(calls, "fourth")
		return nil
	})
	assert.Nil(err)
	assert.Equal([]string{"third", "second", "first", "fourth"}, calls)
}
