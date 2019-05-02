package reverseproxy

import (
	"net/http"
	"time"

	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

// NewHTTPRedirect returns a new HTTPRedirect which redirects HTTP to HTTPS
func NewHTTPRedirect() *HTTPRedirect {
	return &HTTPRedirect{}
}

// HTTPRedirect redirects HTTP to HTTPS
type HTTPRedirect struct {
	Log logger.Log
}

// ServeHTTP redirects HTTP to HTTPS
func (hr *HTTPRedirect) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	w := NewResponseWriter(rw)

	if hr.Log != nil {
		hr.Log.Trigger(req.Context(), logger.NewHTTPRequestEvent(req))

		start := time.Now()
		defer func() {
			wre := logger.NewHTTPResponseEvent(req,
				logger.OptHTTPResponseStatusCode(http.StatusMovedPermanently),
				logger.OptHTTPResponseContentLength(w.ContentLength()),
				logger.OptHTTPResponseElapsed(time.Since(start)),
			)

			if value := rw.Header().Get("Content-Type"); len(value) > 0 {
				wre.ContentType = value
			}
			if value := rw.Header().Get("Content-Encoding"); len(value) > 0 {
				wre.ContentEncoding = value
			}
			hr.Log.Trigger(req.Context(), wre)
		}()
	}

	req.URL.Scheme = webutil.SchemeHTTPS
	if req.URL.Host == "" {
		req.URL.Host = req.Host
	}

	http.Redirect(rw, req, req.URL.String(), http.StatusMovedPermanently)
}
