/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
