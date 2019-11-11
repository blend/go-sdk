package webutil

import (
	"net/http"
	"strings"
)

// HeaderLastValue returns the last value of a potential csv of headers.
func HeaderLastValue(headers http.Header, key string) (string, bool) {
	if headerVal := headers.Get(key); headerVal != "" {
		if !strings.ContainsRune(headerVal, ',') {
			return strings.TrimSpace(headerVal), true
		}
		vals := strings.Split(headerVal, ",")
		return strings.TrimSpace(vals[len(vals)-1]), true
	}
	return "", false
}
