/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package r2

import (
	"net/url"

	"github.com/blend/go-sdk/webutil"
)

// OptPostForm sets the request post form and the content type.
func OptPostForm(postForm url.Values) Option {
	return RequestOption(webutil.OptPostForm(postForm))
}

// OptPostFormValue sets a request post form value.
func OptPostFormValue(key, value string) Option {
	return RequestOption(webutil.OptPostFormValue(key, value))
}
