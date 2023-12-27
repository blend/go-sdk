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
	"net/url"
)

// NewMockRequest creates a mock request.
func NewMockRequest(method, path string) *http.Request {
	return &http.Request{
		Method:		method,
		Proto:		"http",
		ProtoMajor:	1,
		ProtoMinor:	1,
		Host:		"localhost:8080",
		RemoteAddr:	"127.0.0.1:8080",
		RequestURI:	path,
		Header: http.Header{
			HeaderUserAgent: []string{"go-sdk test"},
		},
		URL: &url.URL{
			Scheme:		"http",
			Host:		"localhost",
			Path:		path,
			RawPath:	path,
		},
	}
}

// NewMockRequestWithCookie creates a mock request with a cookie attached to it.
func NewMockRequestWithCookie(method, path, cookieName, cookieValue string) *http.Request {
	req := NewMockRequest(method, path)
	req.AddCookie(&http.Cookie{
		Name:	cookieName,
		Domain:	"localhost",
		Path:	"/",
		Value:	cookieValue,
	})
	return req
}
