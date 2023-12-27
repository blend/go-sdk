/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
