/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
