/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package grpcutil

import (
	"google.golang.org/grpc"
)

// StreamServerChain reads the middleware variadic args and organizes the calls recursively in the order they appear.
func StreamServerChain(interceptors ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	// if we don't have interceptors, return a no-op.
	if len(interceptors) == 0 {
		return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			return handler(srv, ss)
		}
	}

	// if we only have one interceptor, return it
	if len(interceptors) == 1 {
		return interceptors[0]
	}

	// nest the interceptors
	var nest = func(a, b grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
		if b == nil {
			return a
		}
		return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			curried := func(isrv interface{}, istream grpc.ServerStream) error {
				return b(isrv, istream, info, handler)
			}
			return a(srv, ss, info, curried)
		}
	}

	var outer grpc.StreamServerInterceptor
	for _, step := range interceptors {
		outer = nest(step, outer)
	}
	return outer
}
