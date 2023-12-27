/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"io"

	"github.com/blend/go-sdk/webutil"
)

// OptBody sets the post body on the request.
func OptBody(contents io.ReadCloser) Option {
	return RequestOption(webutil.OptBody(contents))
}

// OptBodyBytes sets the post body on the request from a byte array.
func OptBodyBytes(contents []byte) Option {
	return RequestOption(webutil.OptBodyBytes(contents))
}
