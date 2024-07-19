/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"net/http"
)

// OptDisableKeepAlives disables keep alives.
func OptDisableKeepAlives(disableKeepAlives bool) Option {
	return func(r *Request) error {
		if r.Client == nil {
			r.Client = &http.Client{}
		}
		if r.Client.Transport == nil {
			r.Client.Transport = &http.Transport{}
		}
		if typed, ok := r.Client.Transport.(*http.Transport); ok {
			typed.DisableKeepAlives = disableKeepAlives
		}
		return nil
	}
}
