/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import "github.com/blend/go-sdk/webutil"

// OptPostedFiles adds multipart uploads to the request.
//
// Usage note: this option will also encode any currently provided
// post form fields into the body as well, so you should make this the
// last option in a list to capture those fields.
func OptPostedFiles(files ...webutil.PostedFile) Option {
	return RequestOption(webutil.OptPostedFiles(files...))
}
