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

// EnsureHTTPTransport ensures the http client's transport
// is set and that it is an *http.Transport.
//
// It will return an error `ErrInvalidTransport` if it
// is set to something other than *http.Transport.
func EnsureHTTPTransport(r *Request) (*http.Transport, error) {
	if r.Client == nil {
		r.Client = &http.Client{}
	}
	if r.Client.Transport == nil {
		r.Client.Transport = &http.Transport{}
	}
	typed, ok := r.Client.Transport.(*http.Transport)
	if r.Client.Transport != nil && !ok {
		return nil, ex.New(ErrInvalidTransport)
	}
	if typed == nil {
		typed = &http.Transport{}
		r.Client.Transport = typed
	}
	return typed, nil
}
