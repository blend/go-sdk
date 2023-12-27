/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package grpcutil

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// UnaryClientTimeout returns a unary client interceptor.
func UnaryClientTimeout(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		timeoutCtx, done := context.WithTimeout(ctx, timeout)
		defer done()
		return invoker(timeoutCtx, method, req, reply, cc, opts...)
	}
}
