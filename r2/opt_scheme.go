package r2

import (
	"net/url"

	"github.com/blend/go-sdk/ex"
)

// OptScheme sets the url scheme.
func OptScheme(scheme string) Option {
	return func(r *Request) error {
		if r.Request == nil {
			return ex.New(ErrRequestUnset)
		}
		if r.Request.URL == nil {
			r.Request.URL = &url.URL{}
		}
		r.Request.URL.Scheme = scheme
		return nil
	}
}
