/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package r2

import (
	"fmt"
	"net/url"

	"github.com/blend/go-sdk/ex"
)

// OptPort sets a custom port for the request url.
func OptPort(port int32) Option {
	return func(r *Request) error {
		if r.Request == nil {
			return ex.New(ErrRequestUnset)
		}
		if r.Request.URL == nil {
			r.Request.URL = &url.URL{}
		}
		r.Request.URL.Host = fmt.Sprintf("%s:%d", r.Request.URL.Hostname(), port)
		return nil
	}
}
