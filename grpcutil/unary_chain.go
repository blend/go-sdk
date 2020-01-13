package grpcutil

import (
	"context"

	"google.golang.org/grpc"
)

// UnaryChain reads the middleware variadic args and organizes the calls recursively in the order they appear.
func UnaryChain(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	// if we don't have interceptors, return a no-op.
	if len(interceptors) == 0 {
		return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			return handler(ctx, req)
		}
	}
	// if we only have one interceptor, return it
	if len(interceptors) == 1 {
		return interceptors[0]
	}

	// nest the interceptors
	var nest = func(a, b grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
		if b == nil {
			return a
		}
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			curried := func(ictx context.Context, ireq interface{}) (interface{}, error) {
				return b(ictx, ireq, info, handler)
			}
			return a(ctx, req, info, curried)
		}
	}

	var outer grpc.UnaryServerInterceptor
	for _, step := range interceptors {
		outer = nest(step, outer)
	}
	return outer
}
