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

func TestStreamServerChain(t *testing.T) {
	assert := assert.New(t)

	var calls []string
	combined := StreamServerChain(
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
