/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package grpcutil

import (
	"context"

	"google.golang.org/grpc"
)

// StreamClientChain creates a single interceptor out of a chain of many interceptors.
//
// Execution is done in left-to-right order, including passing of context.
// For example ChainStreamClient(one, two, three) will execute one before two before three.
func StreamClientChain(interceptors ...grpc.StreamClientInterceptor) grpc.StreamClientInterceptor {
	n := len(interceptors)

	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		chainer := func(currentInter grpc.StreamClientInterceptor, currentStreamer grpc.Streamer) grpc.Streamer {
			return func(currentCtx context.Context, currentDesc *grpc.StreamDesc, currentConn *grpc.ClientConn, currentMethod string, currentOpts ...grpc.CallOption) (grpc.ClientStream, error) {
				return currentInter(currentCtx, currentDesc, currentConn, currentMethod, currentStreamer, currentOpts...)
			}
		}

		chainedStreamer := streamer
		for i := n - 1; i >= 0; i-- {
			chainedStreamer = chainer(interceptors[i], chainedStreamer)
		}

		return chainedStreamer(ctx, desc, cc, method, opts...)
	}
}
