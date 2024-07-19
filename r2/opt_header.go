/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"net/http"

	"github.com/blend/go-sdk/webutil"
)

// OptHeader sets the request headers.
func OptHeader(headers http.Header) Option {
	return RequestOption(webutil.OptHeader(headers))
}

// OptHeaderValue adds or sets a header value.
func OptHeaderValue(key, value string) Option {
	return RequestOption(webutil.OptHeaderValue(key, value))
}
