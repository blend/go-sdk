/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import "github.com/blend/go-sdk/webutil"

// OptHTTPClientTrace sets the outgoing httptrace.ClientTrace context.
func OptHTTPClientTrace(ht *webutil.HTTPTrace) Option {
	return RequestOption(webutil.OptHTTPClientTrace(ht))
}
