package reverseproxy

import (
	"net/http"

	"github.com/blend/go-sdk/logger"
)

// ProxyOption is an option for the proxy.
type ProxyOption func(*Proxy)

// OptProxyLog sets the logger on the proxy.
func OptProxyLog(log logger.Log) ProxyOption {
	return func(p *Proxy) {
		p.Log = log
	}
}

// OptProxyUpstream adds an upstream to the proxy.
func OptProxyUpstream(upstream *Upstream) ProxyOption {
	return func(p *Proxy) {
		p.Upstreams = append(p.Upstreams, upstream)
	}
}

// OptProxyHeaderValue sets a header value for all requests to the proxy's upstreams.
func OptProxyHeaderValue(key, value string) ProxyOption {
	return func(p *Proxy) {
		if p.Headers == nil {
			p.Headers = http.Header{}
		}
		p.Headers.Set(key, value)
	}
}
