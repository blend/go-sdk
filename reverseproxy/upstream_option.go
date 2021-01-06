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
