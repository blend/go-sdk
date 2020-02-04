package webutil

import "net/http"

// NoFollowRedirects returns an http client redirect checker that returns the
// http.ErrUseLastResponse error.
func NoFollowRedirects() func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
}
