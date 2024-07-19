/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"net/url"

	"github.com/blend/go-sdk/webutil"
)

// OptQuery sets the full querystring.
func OptQuery(query url.Values) Option {
	return RequestOption(webutil.OptQuery(query))
}

// OptQueryValue sets a query value.
func OptQueryValue(key, value string) Option {
	return RequestOption(webutil.OptQueryValue(key, value))
}

// OptQueryValueAdd adds a query value.
func OptQueryValueAdd(key, value string) Option {
	return RequestOption(webutil.OptQueryValueAdd(key, value))
}
