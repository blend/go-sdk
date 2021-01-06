package r2

import (
	"net/http"
)

// OptTransport sets the client transport for a request.
func OptTransport(transport http.RoundTripper) Option {
	return func(r *Request) error {
		if r.Client == nil {
			r.Client = &http.Client{}
		}
		r.Client.Transport = transport
		return nil
	}
}
