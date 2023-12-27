/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
