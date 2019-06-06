package r2

import (
	"fmt"
	"net/url"
)

// OptPath sets the url path.
func OptPath(path string) Option {
	return func(r *Request) error {
		if r.URL == nil {
			r.URL = &url.URL{}
		}
		r.URL.Path = path
		return nil
	}
}

// OptPathf sets the url path based on a format and arguments.
func OptPathf(format string, args ...interface{}) Option {
	return func(r *Request) error {
		if r.URL == nil {
			r.URL = &url.URL{}
		}
		r.URL.Path = fmt.Sprintf(format, args...)
		return nil
	}
}
