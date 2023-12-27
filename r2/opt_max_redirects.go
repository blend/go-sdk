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

	"github.com/blend/go-sdk/ex"
)

// OptMaxRedirects tells the http client to only follow a given
// number of redirects, overriding the standard library default of 10.
// Use the companion helper `ErrIsTooManyRedirects` to test if the returned error
// from a call indicates the redirect limit was reached.
func OptMaxRedirects(maxRedirects int) Option {
	return func(r *Request) error {
		if r.Client == nil {
			r.Client = &http.Client{}
		}
		r.Client.CheckRedirect = func(r *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				// NOTE: if we return *just* http.ErrUseLastResponse here
				// http.Client will short circuit, treating the returned
				// error as a sentinel value and return the last http response with a nil error.
				// We have an explicit function to test for the "exception" form of the error
				// `ErrIsUseLastResponse` that can be used to assert if we returned
				// the error here.
				return ex.New(http.ErrUseLastResponse)
			}
			return nil
		}
		return nil
	}
}
