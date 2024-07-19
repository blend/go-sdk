/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"net/url"

	"github.com/blend/go-sdk/ex"
)

// OptURL sets the url of a request.
func OptURL(rawURL string) Option {
	return func(r *Request) error {
		if r.Request == nil {
			return ex.New(ErrRequestUnset)
		}
		var err error
		r.Request.URL, err = url.Parse(rawURL)
		return ex.New(err)
	}
}
