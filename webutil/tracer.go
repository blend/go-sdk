/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
