/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import "net/http"

// GetUserAgent gets a user agent from a request.
func GetUserAgent(r *http.Request) string {
	return r.UserAgent()
}
