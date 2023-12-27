/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import (
	"net/http"
	"strings"
)

// GetProto gets the request proto.
// X-FORWARDED-PROTO is checked first, then the original request proto is used.
func GetProto(r *http.Request) (scheme string) {
	if r == nil {
		return
	}

	// Retrieve the scheme from X-Forwarded-Proto.
	if proto, ok := HeaderLastValue(r.Header, HeaderXForwardedProto); ok {
		scheme = strings.ToLower(proto)
	} else if proto, ok = HeaderLastValue(r.Header, HeaderXForwardedScheme); ok {
		scheme = strings.ToLower(proto)
	} else if proto, ok = HeaderLastValue(r.Header, HeaderForwarded); ok {
		// match should contain at least two elements if the protocol was
		// specified in the Forwarded header. The first element will always be
		// the 'proto=' capture, which we ignore. In the case of multiple proto
		// parameters (invalid) we only extract the first.
		if match := protoRegex.FindStringSubmatch(proto); len(match) > 1 {
			scheme = strings.ToLower(match[1])
		}
	} else if r.URL != nil {
		scheme = strings.ToLower(r.URL.Scheme)
	}
	return
}
