package reverseproxy

import (
	"fmt"
	"net/http"
	"os"

	"github.com/blend/go-sdk/ex"
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
	Headers          http.Header
	Log              logger.Log
	Upstreams        []*Upstream
	Resolver         Resolver
	Tracer           webutil.HTTPTracer
	TransformRequest TransformRequest
}

// ServeHTTP is the http entrypoint.
func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var err error
	var tf webutil.HTTPTraceFinisher
	defer func() {
		// NOTE: This uses the outer scope's `err` by design. This way updates
		//       to `err` will be reflected on (deferred) exit.
		r := recover()

		// Wrap the error with the reason for the panic.
		if rErr, ok := r.(error); ok {
			err = ex.Nest(err, rErr)
		} else {
			err = ex.Nest(err, ex.New(r))
		}

		// Finish the span, if open
		if tf != nil {
			tf.Finish(err)
		}

		// Log or print the cause of the panic.
		if r == nil {
			return
		}
		if p.Log != nil {
			p.Log.Fatalf("%v", r)
		} else {
			fmt.Fprintf(os.Stderr, "%v\n", r)
		}
	}()

	if p.Tracer != nil {
		tf, req = p.Tracer.Start(req)
	}

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
	if p.TransformRequest != nil {
		err := p.TransformRequest(req)
		if err != nil {
			logger.MaybeError(p.Log, err)
		}
	}

	upstream.ServeHTTP(rw, req)
}
