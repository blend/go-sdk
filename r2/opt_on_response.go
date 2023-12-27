/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package r2

import (
	"net/http"
	"time"
)

// OnResponseListener is an on response listener.
//
// The time.Time is given as the start time of the request in the UTC timezone. To compute the elapsed time
// you would subtract from the current time in UTC i.e. `time.Now().UTC().Sub(startTime)`.
type OnResponseListener func(*http.Request, *http.Response, time.Time, error) error

// OptOnResponse adds an on response listener.
// If an OnResponse listener has already been addded, it will be merged with the existing listener.
func OptOnResponse(listener OnResponseListener) Option {
	return func(r *Request) error {
		r.OnResponse = append(r.OnResponse, listener)
		return nil
	}
}
