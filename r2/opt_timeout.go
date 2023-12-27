/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package r2

import (
	"net/http"
	"time"
)

// OptTimeout sets the client timeout.
func OptTimeout(d time.Duration) Option {
	return func(r *Request) error {
		if r.Client == nil {
			r.Client = &http.Client{}
		}
		r.Client.Timeout = d
		return nil
	}
}
