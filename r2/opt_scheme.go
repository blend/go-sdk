package r2

import (
	"net/http"
	"net/url"
)

// OptScheme sets the url scheme.
func OptScheme(scheme string) Option {
	return func(r *Request) error {
		if r.Request == nil {
			r.Request = &http.Request{}
		}
		if r.URL == nil {
			r.URL = &url.URL{}
		}
		r.URL.Scheme = scheme
		return nil
	}
}
