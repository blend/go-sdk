/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package reverseproxy

import (
	"net/http"

	"github.com/blend/go-sdk/webutil"
)

// HTTPRedirect redirects HTTP to HTTPS
type HTTPRedirect struct {
	RedirectScheme	string
	RedirectHost	string
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
