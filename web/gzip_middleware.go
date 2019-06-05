package web

import (
	"github.com/blend/go-sdk/webutil"
)

// GZip is a middleware the implements gzip compression for requests that opt into it.
func GZip(action Action) Action {
	return func(r *Ctx) Result {
		w := r.Response
		if webutil.HeaderAny(r.Request.Header, HeaderAcceptEncoding, ContentEncodingGZIP) {
			w.Header().Set(HeaderContentEncoding, ContentEncodingGZIP)
			w.Header().Set(HeaderVary, "Accept-Encoding")
			r.Response = webutil.NewGZipResponseWriter(w)
		}
		return action(r)
	}
}
