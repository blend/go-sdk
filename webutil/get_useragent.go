/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import "net/http"

// GetUserAgent gets a user agent from a request.
func GetUserAgent(r *http.Request) string {
	return r.UserAgent()
}
