/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package grpcutil

import (
	"context"
	"testing"

	"google.golang.org/grpc"

	"github.com/blend/go-sdk/assert"
)

func TestUnaryServerChain(t *testing.T) {
	assert := assert.New(t)

	var calls []string
	combined := UnaryServerChain(
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
