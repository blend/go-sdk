package r2

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/blend/go-sdk/ex"
)

// OptPath sets the url path.
func OptPath(path string) Option {
	return func(r *Request) error {
		if r.Request == nil {
			return ex.New(ErrRequestUnset)
		}
		if r.Request.URL == nil {
			r.Request.URL = &url.URL{}
		}
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		r.Request.URL.Path = path
		return nil
	}
}

// OptPathf sets the url path based on a format and arguments.
func OptPathf(format string, args ...interface{}) Option {
	return OptPath(fmt.Sprintf(format, args...))
}
