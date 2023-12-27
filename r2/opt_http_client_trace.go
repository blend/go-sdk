/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package r2

import "github.com/blend/go-sdk/webutil"

// OptHTTPClientTrace sets the outgoing httptrace.ClientTrace context.
func OptHTTPClientTrace(ht *webutil.HTTPTrace) Option {
	return RequestOption(webutil.OptHTTPClientTrace(ht))
}
