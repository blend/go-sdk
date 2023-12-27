/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"net/http"
	"time"
)

// EventOption is an event option.
type EventOption func(e *Event)

// OptEventRequest sets the response.
func OptEventRequest(req *http.Request) EventOption {
	return func(e *Event) {
		e.Request = req
	}
}

// OptEventResponse sets the response.
func OptEventResponse(res *http.Response) EventOption {
	return func(e *Event) {
		e.Response = res
	}
}

// OptEventElapsed sets the elapsed time.
func OptEventElapsed(elapsed time.Duration) EventOption {
	return func(e *Event) {
		e.Elapsed = elapsed
	}
}

// OptEventBody sets the body.
func OptEventBody(body []byte) EventOption {
	return func(e *Event) {
		e.Body = body
	}
}
