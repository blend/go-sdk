/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

import (
	"github.com/blend/go-sdk/webutil"
)

// GZip is a middleware the implements gzip compression for requests that opt into it.
func GZip(action Action) Action {
	return func(r *Ctx) Result {
		if webutil.HeaderAny(r.Request.Header, webutil.HeaderAcceptEncoding, webutil.ContentEncodingGZIP) {
			r.Response.Header().Set(webutil.HeaderContentEncoding, webutil.ContentEncodingGZIP)
			r.Response.Header().Set(webutil.HeaderVary, webutil.HeaderAcceptEncoding)
			r.Response = webutil.NewGZipResponseWriter(r.Response)
		}
		return action(r)
	}
}
