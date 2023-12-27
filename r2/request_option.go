/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package r2

import "net/http"

// RequestOption translates a webutil.RequestOption to a r2.Option.
func RequestOption(opt func(*http.Request) error) Option {
	return func(r *Request) error {
		return opt(r.Request)
	}
}
