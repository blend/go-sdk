/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package reverseproxy

import (
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/http/httpguts"
)

// MustParseURL parses a url and panics if it's bad.
func MustParseURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return u
}

// RequestCopy does a shallow copy of a request.
func RequestCopy(req *http.Request) *http.Request {
	outreq := new(http.Request)
	*outreq = *req	// includes shallow copies of maps, but okay
	if req.ContentLength == 0 {
		outreq.Body = nil
	}
	return outreq
}

// UpgradeType returns the connection upgrade type.
// This is used by websockt support.
func UpgradeType(h http.Header) string {
	if !httpguts.HeaderValuesContainsToken(h["Connection"], "Upgrade") {
		return ""
	}
	return strings.ToLower(h.Get("Upgrade"))
}
