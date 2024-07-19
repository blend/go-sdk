/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import "github.com/blend/go-sdk/webutil"

// OptBasicAuth is an option that sets the http basic auth.
func OptBasicAuth(username, password string) Option {
	return RequestOption(webutil.OptBasicAuth(username, password))
}
