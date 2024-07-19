/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import "net/http"

// Handler is the most basic route handler.
type Handler func(http.ResponseWriter, *http.Request, *Route, RouteParameters)

// WrapHandler wraps an http.Handler as a Handler.
func WrapHandler(handler http.Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request, _ *Route, _ RouteParameters) {
		handler.ServeHTTP(w, r)
	}
}
