package proxy

import (
	"fmt"
	"net/http"
	"os"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

const (
	// FlagProxyRequest is a logger flag.
	FlagProxyRequest logger.Flag = "proxy.request"
)

// New returns a new proxy.
func New() *Proxy {
	return &Proxy{}
}

// Proxy is a factory for a simple reverse proxy.
type Proxy struct {
	upstreamHeaders http.Header
	log             *logger.Logger
	upstreams       []*Upstream
	resolver        Resolver
}

// WithUpstreamHeaders sets headers to be added to all upstream requests.
// Note: this will overwrite any existing headers.
func (p *Proxy) WithUpstreamHeaders(headers http.Header) *Proxy {
	p.upstreamHeaders = headers
	return p
}

// WithUpstreamHeader adds a single upstream header.
func (p *Proxy) WithUpstreamHeader(key, value string) *Proxy {
	if p.upstreamHeaders == nil {
		p.upstreamHeaders = http.Header{}
	}
	p.upstreamHeaders.Set(key, value)
	return p
}

// UpstreamHeaders returns the upstream headers to add to all upstream requests.
func (p *Proxy) UpstreamHeaders() http.Header {
	return p.upstreamHeaders
}

// WithLogger sets a property and returns the proxy reference.
func (p *Proxy) WithLogger(log *logger.Logger) *Proxy {
	p.log = log
	return p
}

// WithUpstream adds an upstream by URL.
func (p *Proxy) WithUpstream(upstream *Upstream) *Proxy {
	p.upstreams = append(p.upstreams, upstream)
	return p
}

// WithResolver sets a property and returns the proxy reference.
func (p *Proxy) WithResolver(resolver Resolver) *Proxy {
	p.resolver = resolver
	return p
}

// ServeHTTP is the http entrypoint.
func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			if p.log != nil {
				p.log.Fatalf("%v", r)
			} else {
				fmt.Fprintf(os.Stderr, "%v\n", r)
			}
		}
	}()

	// set the default resolver if unset.
	if p.resolver == nil {
		p.resolver = RoundRobinResolver(p.upstreams)
	}

	upstream, err := p.resolver(req, p.upstreams)

	if err != nil {
		p.log.Error(err)
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	if upstream == nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	// Add extra forwarded headers.
	// these are required for a majority of services to function correctly behind
	// a reverse proxy.
	// They are "Add" vs. "Set" in case there are existing values.
	if port := webutil.GetPort(req); port != "" {
		req.Header.Add("X-Forwarded-Port", port)
	}
	if proto := webutil.GetProto(req); proto != "" {
		req.Header.Add("X-Forwarded-Proto", proto)
	}
	// add upstream headers.
	for key, values := range p.upstreamHeaders {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	upstream.ServeHTTP(rw, req)
}
