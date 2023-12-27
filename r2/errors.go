/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"net/http"
	"net/url"

	"github.com/blend/go-sdk/ex"
)

// Error Constants
const (
	// ErrRequestUnset is an error returned from options if they are called on a r2.Request that has not
	// been created by r2.New(), and as a result the underlying request is uninitialized.
	ErrRequestUnset	ex.Class	= "r2; cannot modify request, underlying request unset. please use r2.New()"
	// ErrInvalidTransport is an error returned from options if they are called on a request that has had
	// the transport set to something other than an *http.Transport; this precludes using http.Transport
	// specific options like tls.Config mutators.
	ErrInvalidTransport	ex.Class	= "r2; cannot modify transport, is not *http.Transport"
	// ErrNoContentJSON is returns from sending requests when a no-content status is returned.
	ErrNoContentJSON	ex.Class	= "server returned an http 204 for a request expecting json"
	// ErrNoContentXML is returns from sending requests when a no-content status is returned.
	ErrNoContentXML	ex.Class	= "server returned an http 204 for a request expecting xml"
	// ErrInvalidMethod is an error that is returned from `r2.Request.Do()` if a method
	// is specified on the request that violates the valid charset for HTTP methods.
	ErrInvalidMethod	ex.Class	= "r2; invalid http method"
	// ErrMismatchedPathParameters is an error that is returned from `OptPathParameterized()` if
	// the parameterized path string has a different number of parameters than what was passed as
	// variadic arguments.
	ErrMismatchedPathParameters	ex.Class	= "r2; route parameters provided don't match parameters needed in path"
)

// ErrIsTooManyRedirects returns if the error is too many redirects.
func ErrIsTooManyRedirects(err error) bool {
	if ex.Is(err, http.ErrUseLastResponse) {
		return true
	}
	if typed, ok := err.(*url.Error); ok {
		return ex.Is(typed.Err, http.ErrUseLastResponse)
	}
	return false
}
