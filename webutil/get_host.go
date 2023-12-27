/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import (
	"net/http"
	"strings"
)

// GetHost returns the request host, omitting the port if specified.
func GetHost(r *http.Request) string {
	if r == nil {
		return ""
	}
	tryHeader := func(key string) (string, bool) {
		return HeaderLastValue(r.Header, key)
	}
	for _, header := range []string{HeaderXForwardedHost} {
		if headerVal, ok := tryHeader(header); ok {
			return headerVal
		}
	}
	if r.URL != nil && len(r.URL.Host) > 0 {
		return r.URL.Host
	}
	if strings.Contains(r.Host, ":") {
		return strings.SplitN(r.Host, ":", 2)[0]
	}
	return r.Host
}

// GetHostStrict returns the request host, omitting the port if specified,
// and does not consider `X-Forwarded-Host` headers.
func GetHostStrict(r *http.Request) string {
	if r == nil {
		return ""
	}
	if r.Host != "" {
		if strings.Contains(r.Host, ":") {
			return strings.SplitN(r.Host, ":", 2)[0]
		}
		return r.Host
	}
	return ""
}
