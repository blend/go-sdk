/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"net/http"
)

// OptTransport sets the client transport for a request.
func OptTransport(transport http.RoundTripper) Option {
	return func(r *Request) error {
		if r.Client == nil {
			r.Client = &http.Client{}
		}
		r.Client.Transport = transport
		return nil
	}
}
