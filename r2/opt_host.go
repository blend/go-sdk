/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"net/url"

	"github.com/blend/go-sdk/ex"
)

// OptHost sets the url host.
func OptHost(host string) Option {
	return func(r *Request) error {
		if r.Request == nil {
			return ex.New(ErrRequestUnset)
		}
		if r.Request.URL == nil {
			r.Request.URL = &url.URL{}
		}
		r.Request.URL.Host = host
		return nil
	}
}

// OptServiceHostName sets the service hostName within an outgoing request context.
// Service host names are useful for aggregating outbound service requests where
// host headers should be deduped instead.
func OptServiceHostName(serviceHostName string) Option {
	return func(r *Request) error {
		if r.Request == nil {
			return ex.New(ErrRequestUnset)
		}

		ctx := r.Request.Context()
		r.WithContext(WithServiceHostName(ctx, serviceHostName))

		return nil
	}
}
