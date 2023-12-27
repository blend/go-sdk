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
	"fmt"
	"net/url"
	"strings"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/stringutil"
)

// OptPath sets the url path.
func OptPath(path string) Option {
	return func(r *Request) error {
		if r.Request == nil {
			return ex.New(ErrRequestUnset)
		}
		if r.Request.URL == nil {
			r.Request.URL = &url.URL{}
		}
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		r.Request.URL.Path = path
		return nil
	}
}

// OptPathf sets the url path based on a format and arguments.
func OptPathf(format string, args ...interface{}) Option {
	return OptPath(fmt.Sprintf(format, args...))
}

// OptPathParameterized sets the url path based on a parameterized path and arguments.
// Parameterized paths should appear in the same format as paths you would add to your
// web app (ex. `/resource/:resource_id`).
func OptPathParameterized(format string, params map[string]string) Option {
	return func(r *Request) error {
		if r.Request == nil {
			return ex.New(ErrRequestUnset)
		}

		if !strings.HasPrefix(format, "/") {
			format = "/" + format
		}

		path, err := stringutil.ReplacePathParameters(format, params)
		if err != nil {
			return err
		}

		ctx := r.Request.Context()
		r.WithContext(WithParameterizedPath(ctx, format))

		return OptPath(path)(r)
	}
}
