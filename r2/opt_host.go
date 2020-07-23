package r2

import (
	"net/url"

	"github.com/blend/go-sdk/ex"
)

// OptHost sets the url host.
func OptHost(host string) Option {
	return func(r *Request) error {
		if r.Request == nil {
			return ex.New(ErrRequestUnset)
		}
		if r.Request.URL == nil {
			r.Request.URL = &url.URL{}
		}
		r.Request.URL.Host = host
		return nil
	}
}
