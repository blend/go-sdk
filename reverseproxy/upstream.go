package reverseproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/blend/go-sdk/logger"
)

// NewUpstream returns a new upstram.
func NewUpstream(target *url.URL) *Upstream {
	return &Upstream{
		URL:          target,
		ReverseProxy: httputil.NewSingleHostReverseProxy(target),
	}
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

// ServeHTTP
func (u *Upstream) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	w := NewResponseWriter(rw)

	if u.Log != nil {
		u.Log.Trigger(req.Context(), logger.NewHTTPRequestEvent(req))

		start := time.Now()
		defer func() {
			wre := logger.NewHTTPResponseEvent(req,
				logger.OptHTTPResponseStatusCode(w.StatusCode()),
				logger.OptHTTPResponseContentLength(w.ContentLength()),
				logger.OptHTTPResponseElapsed(time.Since(start)),
			)

			if value := w.Header().Get("Content-Type"); len(value) > 0 {
				wre.ContentType = value
			}
			if value := w.Header().Get("Content-Encoding"); len(value) > 0 {
				wre.ContentEncoding = value
			}

			u.Log.Trigger(req.Context(), wre)
		}()
	}

	// Add extra forwarded headers.
	// these are required for a majority of services to function correctly behind
	// a reverse proxy.
	w.Header().Set("X-Forwarded-Port", req.URL.Port())
	w.Header().Set("X-Forwarded-Proto", req.URL.Scheme)

	u.ReverseProxy.ServeHTTP(w, req)
}
