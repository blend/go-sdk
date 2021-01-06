package reverseproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"golang.org/x/net/http2"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

// NewUpstream returns a new upstram.
func NewUpstream(target *url.URL) *Upstream {
	rp := httputil.NewSingleHostReverseProxy(target)
	u := &Upstream{
		URL:          target,
		ReverseProxy: rp,
	}
	// NOTE: This creates a reference cycle `u -> rp -> u`.
	rp.ErrorHandler = u.errorHandler
	return u
}

// Upstream represents a proxyable server.
type Upstream struct {
	// Name is the name of the upstream.
	Name string
	// Log is a logger agent.
	Log logger.Log
	// URL represents the target of the upstream.
	URL *url.URL
	// ReverseProxy is what actually forwards requests.
	ReverseProxy *httputil.ReverseProxy
}

// UseHTTP2 sets the upstream to use http2.
func (u *Upstream) UseHTTP2() error {
	if u.ReverseProxy.Transport == nil {
		u.ReverseProxy.Transport = &http.Transport{}
	}
	if typed, ok := u.ReverseProxy.Transport.(*http.Transport); ok {
		if err := http2.ConfigureTransport(typed); err != nil {
			return ex.New(err)
		}
	}
	return nil
}

// ServeHTTP
func (u *Upstream) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	w := webutil.NewResponseWriter(rw)

	if u.Log != nil {
		u.Log.TriggerContext(req.Context(), webutil.NewHTTPRequestEvent(req))

		start := time.Now()
		defer func() {
			wre := webutil.NewHTTPResponseEvent(req,
				webutil.OptHTTPResponseStatusCode(w.StatusCode()),
				webutil.OptHTTPResponseContentLength(w.ContentLength()),
				webutil.OptHTTPResponseElapsed(time.Since(start)),
			)

			if value := w.Header().Get("Content-Type"); len(value) > 0 {
				wre.ContentType = value
			}
			if value := w.Header().Get("Content-Encoding"); len(value) > 0 {
				wre.ContentEncoding = value
			}

			u.Log.TriggerContext(req.Context(), wre)
		}()
	}

	u.ReverseProxy.ServeHTTP(w, req)
}

// errorHandler is intended to be used as an `(net/http/httputil).ReverseProxy.ErrorHandler`
// This implementation is based on:
// https://github.com/golang/go/blob/go1.13.6/src/net/http/httputil/reverseproxy.go#L151-L154
func (u *Upstream) errorHandler(rw http.ResponseWriter, req *http.Request, err error) {
	logger.MaybeErrorfContext(req.Context(), u.Log, "http: proxy error: %v", err)
	rw.WriteHeader(http.StatusBadGateway)
}
