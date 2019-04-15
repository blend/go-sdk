package r2

import "net/url"

// OptScheme sets the url scheme.
func OptScheme(scheme string) Option {
	return func(r *Request) error {
		if r.URL == nil {
			r.URL = &url.URL{}
		}
		r.URL.Scheme = scheme
		return nil
	}
}
