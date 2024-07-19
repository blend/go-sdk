/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package reverseproxy

import (
	"net"
	"net/http"

	"github.com/blend/go-sdk/webutil"
)

// UpstreamOption sets upstream options.
type UpstreamOption func(*Upstream)

// OptUpstreamDial sets the dial options for the upstream.
func OptUpstreamDial(opts ...webutil.DialOption) UpstreamOption {
	return func(u *Upstream) {
		if u.ReverseProxy.Transport == nil {
			u.ReverseProxy.Transport = new(http.Transport)
		}
		if typed, ok := u.ReverseProxy.Transport.(*http.Transport); ok {
			dialer := new(net.Dialer)
			for _, opt := range opts {
				opt(dialer)
			}
			typed.DialContext = dialer.DialContext
		}
	}
}

// OptUpstreamModifyResponse sets the dial options for the upstream.
func OptUpstreamModifyResponse(modifyResponse func(*http.Response) error) UpstreamOption {
	return func(u *Upstream) {
		u.ReverseProxy.ModifyResponse = modifyResponse
	}
}
