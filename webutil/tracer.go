/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"net/http"
)

// HTTPTracer is a simplified version of `Tracer` intended for a raw
// `(net/http).Request`. It returns a "new" request the request context may
// be modified after opening a span.
type HTTPTracer interface {
	Start(*http.Request) (HTTPTraceFinisher, *http.Request)
}

// HTTPTraceFinisher is a simplified version of `TraceFinisher`.
type HTTPTraceFinisher interface {
	Finish(int, error)
}
