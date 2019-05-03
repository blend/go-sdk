package r2

import (
	"net/http"
	"net/url"
)

// OptPath sets the url path.
func OptPath(path string) Option {
	return func(r *Request) error {
		if r.Request == nil {
			r.Request = &http.Request{}
		}
		if r.URL == nil {
			r.URL = &url.URL{}
		}
		r.URL.Path = path
		return nil
	}
}
