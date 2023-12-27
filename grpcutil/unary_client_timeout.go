/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
