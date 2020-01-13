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
		u.Log.Trigger(req.Context(), webutil.NewHTTPRequestEvent(req))

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

			u.Log.Trigger(req.Context(), wre)
		}()
	}

	u.ReverseProxy.ServeHTTP(w, req)
}
