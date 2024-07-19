/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import "github.com/blend/go-sdk/webutil"

// OptMethod sets the request method.
func OptMethod(method string) Option {
	return RequestOption(webutil.OptMethod(method))
}

// OptGet sets the request method.
func OptGet() Option {
	return RequestOption(webutil.OptGet())
}

// OptPost sets the request method.
func OptPost() Option {
	return RequestOption(webutil.OptPost())
}

// OptPut sets the request method.
func OptPut() Option {
	return RequestOption(webutil.OptPut())
}

// OptPatch sets the request method.
func OptPatch() Option {
	return RequestOption(webutil.OptPatch())
}

// OptDelete sets the request method.
func OptDelete() Option {
	return RequestOption(webutil.OptDelete())
}
