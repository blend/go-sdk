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

// OptClient sets the underlying client on the request.
//
// It is specifically useful to prevent churning allocations on
// sending repeated requests.
func OptClient(client *http.Client) Option {
	return func(r *Request) error {
		r.Client = client
		return nil
	}
}
