/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import "net/http"

// RequestOption translates a webutil.RequestOption to a r2.Option.
func RequestOption(opt func(*http.Request) error) Option {
	return func(r *Request) error {
		return opt(r.Request)
	}
}
