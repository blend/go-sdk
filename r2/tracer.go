/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
