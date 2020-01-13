package grpcutil

import (
	"context"
	"testing"
	"google.golang.org/grpc"
	"github.com/blend/go-sdk/assert"
)

func TestUnaryChain(t *testing.T) {
	assert := assert.New(t)

	var calls []string
	combined := UnaryChain(
		func(ctx context.Context, args interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			calls = append(calls, "first")
			return handler(ctx, args)
		},
		func(ctx context.Context, args interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			calls = append(calls, "second")
			return handler(ctx, args)
		},
		func(ctx context.Context, args interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			calls = append(calls, "third")
			return handler(ctx, args)
		},
	)

	res, err := combined(context.Background(), nil, nil, func(ctx context.Context, args interface{}) (interface{}, error) {
		calls = append(calls, "fourth")
		return "ok!", nil
	})
	assert.Nil(err)
	assert.Equal("ok!", res)
	assert.Equal([]string{"third", "second", "first", "fourth"}, calls)
}
