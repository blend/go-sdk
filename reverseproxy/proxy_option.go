package reverseproxy

import (
	"net/http"

	"github.com/blend/go-sdk/logger"
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
