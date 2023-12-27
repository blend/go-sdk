/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import (
	"net/http"
)

var (
	_ http.HandlerFunc = HTTPSRedirectFunc
)

// HTTPSRedirectFunc redirects HTTP to HTTPS as an http.HandlerFunc.
func HTTPSRedirectFunc(rw http.ResponseWriter, req *http.Request) {
	req.URL.Scheme = SchemeHTTPS
	if req.URL.Host == "" {
		req.URL.Host = req.Host
	}
	http.Redirect(rw, req, req.URL.String(), http.StatusMovedPermanently)
}
