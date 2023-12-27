/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

import (
	"net/http"

	"github.com/blend/go-sdk/webutil"
)

// BaseHeaders are the default headers added by go-web.
func BaseHeaders() http.Header {
	return http.Header{
		webutil.HeaderServer: []string{PackageName},
	}
}
