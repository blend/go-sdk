package webutil

import (
	"net"
	"net/http"
)

// GetRemoteAddr gets the origin/client ip for a request.
// X-FORWARDED-FOR is checked. If multiple IPs are included the first one is returned
// X-REAL-IP is checked. If multiple IPs are included the first one is returned
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
