/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"net/http"
	"time"
)

// Tracer is a tracer for requests.
type Tracer interface {
	Start(*http.Request) TraceFinisher
}

// TraceFinisher is a finisher for traces.
type TraceFinisher interface {
	Finish(*http.Request, *http.Response, time.Time, error)
}
