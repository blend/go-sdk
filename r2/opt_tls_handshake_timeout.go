/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"time"
)

// OptTLSHandshakeTimeout sets the client transport TLSHandshakeTimeout.
func OptTLSHandshakeTimeout(d time.Duration) Option {
	return func(r *Request) error {
		transport, err := EnsureHTTPTransport(r)
		if err != nil {
			return err
		}
		transport.TLSHandshakeTimeout = d
		return nil
	}
}
