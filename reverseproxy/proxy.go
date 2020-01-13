package reverseproxy

import (
	"fmt"
	"net/http"
	"os"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

const (
	// FlagProxyRequest is a logger flag.
	FlagProxyRequest = "proxy.request"
)

// NewProxy returns a new proxy.
func NewProxy(opts ...ProxyOption) *Proxy {
	p := Proxy{
		Headers: http.Header{},
	}
	for _, opt := range opts {
		opt(&p)
	}
	return &p
}

// Proxy is a factory for a simple reverse proxy.
type Proxy struct {
	Headers   http.Header
	Log       logger.Log
	Upstreams []*Upstream
	Resolver  Resolver
}

// ServeHTTP is the http entrypoint.
func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			if p.Log != nil {
				p.Log.Fatalf("%v", r)
			} else {
				fmt.Fprintf(os.Stderr, "%v\n", r)
			}
		}
	}()

	// set the default resolver if unset.
	if p.Resolver == nil {
		p.Resolver = RoundRobinResolver(p.Upstreams)
	}

	upstream, err := p.Resolver(req, p.Upstreams)

	if err != nil {
		logger.MaybeError(p.Log, err)
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
	for key, values := range p.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	upstream.ServeHTTP(rw, req)
}
