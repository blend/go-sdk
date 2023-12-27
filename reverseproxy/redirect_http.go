/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package reverseproxy

import (
	"net/http"

	"github.com/blend/go-sdk/webutil"
)

// HTTPRedirect redirects HTTP to HTTPS
type HTTPRedirect struct {
	RedirectScheme string
	RedirectHost   string
}

// ServeHTTP redirects HTTP to HTTPS
func (hr HTTPRedirect) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if hr.RedirectScheme != "" {
		req.URL.Scheme = hr.RedirectScheme
	} else {
		req.URL.Scheme = webutil.SchemeHTTPS
	}
	if hr.RedirectHost != "" {
		req.URL.Host = hr.RedirectHost
	}
	if req.URL.Host == "" {
		req.URL.Host = req.Host
	}
	http.Redirect(rw, req, req.URL.String(), http.StatusMovedPermanently)
}
