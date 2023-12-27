/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package reverseproxy

import (
	"net/http"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

// ProxyOption is a function that mutates a proxy.
type ProxyOption func(*Proxy) error

// OptProxyLog sets the proxy logger, as well
// as the logger on any upstreams that are configured.
func OptProxyLog(log logger.Log) ProxyOption {
	return func(p *Proxy) error {
		p.Log = log
		for _, us := range p.Upstreams {
			us.Log = log
		}
		return nil
	}
}

// OptProxyResolver sets the proxy resolver.
func OptProxyResolver(resolver Resolver) ProxyOption {
	return func(p *Proxy) error {
		p.Resolver = resolver
		return nil
	}
}

// OptProxyUpstream adds a proxy upstream.
func OptProxyUpstream(upstream *Upstream) ProxyOption {
	return func(p *Proxy) error {
		p.Upstreams = append(p.Upstreams, upstream)
		return nil
	}
}

// OptProxyAddHeaderValue adds a proxy upstream.
func OptProxyAddHeaderValue(key, value string) ProxyOption {
	return func(p *Proxy) error {
		if p.Headers == nil {
			p.Headers = http.Header{}
		}
		p.Headers.Add(key, value)
		return nil
	}
}

// OptProxySetHeaderValue adds a proxy upstream.
func OptProxySetHeaderValue(key, value string) ProxyOption {
	return func(p *Proxy) error {
		if p.Headers == nil {
			p.Headers = http.Header{}
		}
		p.Headers.Set(key, value)
		return nil
	}
}

// OptProxyDeleteHeader adds a proxy upstream.
func OptProxyDeleteHeader(key string) ProxyOption {
	return func(p *Proxy) error {
		if p.Headers == nil {
			p.Headers = http.Header{}
		}
		p.Headers.Del(key)
		return nil
	}
}

// OptProxyTracer adds a proxy tracer.
func OptProxyTracer(tracer webutil.HTTPTracer) ProxyOption {
	return func(p *Proxy) error {
		p.Tracer = tracer
		return nil
	}
}

// OptProxyTransformRequest sets the `TransformRequest` on a `Proxy`.
func OptProxyTransformRequest(tr TransformRequest) ProxyOption {
	return func(p *Proxy) error {
		p.TransformRequest = tr
		return nil
	}
}
