/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
