package web

import (
	"net/http"
	"strings"

	"github.com/blend/go-sdk/webutil"
)

// GZip is a middleware the implements gzip compression for requests that opt into it.
func GZip() Middleware {
	return func(action Action) Action {
		return func(r *Ctx) Result {
			w := r.Response
			if strings.Contains(r.Request.Header.Get(HeaderAcceptEncoding), ContentEncodingGZIP) {
				w.Header().Set(http.CanonicalHeaderKey(HeaderContentEncoding), ContentEncodingGZIP)
				w.Header().Set("Vary", "Accept-Encoding")
				r.Response = webutil.NewGZipResponseWriter(w)
			}
			return action(r)
		}
	}
}
