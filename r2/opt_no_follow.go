/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import "net/http"

// OptNoFollow tells the http client to not follow redirects returned
// by the remote server.
func OptNoFollow() Option {
	return func(r *Request) error {
		if r.Client == nil {
			r.Client = &http.Client{}
		}
		r.Client.CheckRedirect = func(_ *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
		return nil
	}
}
