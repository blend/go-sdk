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

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

type peerInfoCommonNameKey struct{}

// WithClientCommonName adds a client common name to a context as a value.
// This value will supercede parsing the value.
func WithClientCommonName(ctx context.Context, commonName string) context.Context {
	return context.WithValue(ctx, peerInfoCommonNameKey{}, commonName)
}

// GetClientCommonName fetches the client common name from the context.
func GetClientCommonName(ctx context.Context) (clientCommonName string) {
	if typed, ok := ctx.Value(peerInfoCommonNameKey{}).(string); ok {
		return typed
	}
	if peer, ok := peer.FromContext(ctx); ok {
		if tlsInfo, ok := peer.AuthInfo.(credentials.TLSInfo); ok {
			if len(tlsInfo.State.VerifiedChains) > 0 && len(tlsInfo.State.VerifiedChains[0]) > 0 {
				clientCommonName = tlsInfo.State.VerifiedChains[0][0].Subject.CommonName
			}
		}
	}
	return
}
