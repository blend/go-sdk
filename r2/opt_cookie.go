/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package r2

import (
	"net/http"

	"github.com/blend/go-sdk/webutil"
)

// OptCookie adds a cookie.
func OptCookie(cookie *http.Cookie) Option {
	return RequestOption(webutil.OptCookie(cookie))
}

// OptCookieValue adds a cookie with a given name and value.
func OptCookieValue(name, value string) Option {
	return RequestOption(webutil.OptCookieValue(name, value))
}
