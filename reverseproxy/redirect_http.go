/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package reverseproxy

import (
	"net/http"

	"github.com/blend/go-sdk/webutil"
)

// HTTPRedirect redirects HTTP to HTTPS
type HTTPRedirect struct{}

// ServeHTTP redirects HTTP to HTTPS
func (hr HTTPRedirect) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	req.URL.Scheme = webutil.SchemeHTTPS
	if req.URL.Host == "" {
		req.URL.Host = req.Host
	}

	http.Redirect(rw, req, req.URL.String(), http.StatusMovedPermanently)
}
