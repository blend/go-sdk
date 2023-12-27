/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import (
	"net/http"
)

// GetPort returns the port for a given request.
func GetPort(r *http.Request) string {
	if r == nil {
		return ""
	}

	tryHeader := func(key string) (string, bool) {
		return HeaderLastValue(r.Header, key)
	}
	for _, header := range []string{HeaderXForwardedPort} {
		if headerVal, ok := tryHeader(header); ok {
			return headerVal
		}
	}
	return ""
}
