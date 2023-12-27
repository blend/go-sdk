/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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
