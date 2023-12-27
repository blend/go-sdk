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

// OnRequestListener is an a listener for on request events.
type OnRequestListener func(*http.Request) error

// OptOnRequest sets an on request listener.
func OptOnRequest(listener OnRequestListener) Option {
	return func(r *Request) error {
		r.OnRequest = append(r.OnRequest, listener)
		return nil
	}
}
