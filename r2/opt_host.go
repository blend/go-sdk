package r2

import (
	"net/http"
	"net/url"
)

// OptHost sets the url host.
func OptHost(host string) Option {
	return func(r *Request) error {
		if r.Request == nil {
			r.Request = &http.Request{}
		}
		if r.URL == nil {
			r.URL = &url.URL{}
		}
		r.URL.Host = host
		return nil
	}
}
