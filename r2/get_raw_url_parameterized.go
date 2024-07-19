/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"fmt"
	"net/http"
	"strings"
)

// GetRawURLParameterized gets a URL string with named route parameters in place of
// the raw path for a request. Useful for outbound request aggregation for
// metrics and tracing when route parameters are involved.
// Relies on the request's context storing the optional hostname and/or parameterized path,
// otherwise will default to returning the request `URL`'s `String()`.
func GetRawURLParameterized(req *http.Request) string {
	if req == nil || req.URL == nil {
		return ""
	}
	url := req.URL

	hostName := GetServiceHostName(req.Context())
	path := GetParameterizedPath(req.Context())
	if path == "" && hostName == "" {
		return url.String()
	}
	if hostName == "" {
		hostName = url.Host
	} else {
		// Using similar formatting as DD to signal this is a parameterized value
		// https://docs.datadoghq.com/tracing/troubleshooting/quantization/#overview
		hostName = fmt.Sprintf("{%s}", hostName)
	}
	if path == "" {
		path = url.Path
	}

	// Stripped down version of "net/url" `URL.String()`
	var buf strings.Builder
	if url.Scheme != "" {
		buf.WriteString(url.Scheme)
		buf.WriteByte(':')
	}
	if url.Scheme != "" || hostName != "" {
		if hostName != "" || path != "" {
			buf.WriteString("//")
		}
		if hostName != "" {
			buf.WriteString(hostName)
		}
	}
	if !strings.HasPrefix(path, "/") {
		buf.WriteByte('/')
	}
	buf.WriteString(path)
	return buf.String()
}
