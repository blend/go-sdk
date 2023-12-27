/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

// github:codeowner @blend/seceng

import (
	"net"
	"net/http"
)

// GetRemoteAddr gets the origin/client ip for a request.
// X-FORWARDED-FOR is checked. If multiple IPs are included the first one is returned
// X-REAL-IP is checked. If multiple IPs are included the last one is returned
// Finally r.RemoteAddr is used
// Only benevolent services will allow access to the real IP.
func GetRemoteAddr(r *http.Request) string {
	if r == nil {
		return ""
	}
	tryHeader := func(key string) (string, bool) {
		return HeaderLastValue(r.Header, key)
	}
	for _, header := range []string{HeaderXForwardedFor, HeaderXRealIP} {
		if headerVal, ok := tryHeader(header); ok {
			return headerVal
		}
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
