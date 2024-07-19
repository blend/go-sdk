/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"net/http"
	"time"

	"github.com/blend/go-sdk/logger"
)

// HTTPLogged returns a middleware that logs a request.
func HTTPLogged(log logger.Triggerable) Middleware {
	return func(action http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			start := time.Now()
			w := NewStatusResponseWriter(rw)
			defer func() {
				responseEvent := NewHTTPRequestEvent(req,
					OptHTTPRequestStatusCode(w.StatusCode()),
					OptHTTPRequestContentLength(w.ContentLength()),
					OptHTTPRequestElapsed(time.Since(start)),
				)
				if w.Header() != nil {
					responseEvent.ContentType = w.Header().Get(HeaderContentType)
					responseEvent.ContentEncoding = w.Header().Get(HeaderContentEncoding)
				}
				logger.MaybeTriggerContext(
					req.Context(),
					log,
					responseEvent,
				)
			}()
			action(w, req)
		}
	}
}
