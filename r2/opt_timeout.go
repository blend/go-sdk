/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
