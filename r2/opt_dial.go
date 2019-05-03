package r2

import (
	"net"
	"net/http"
)

// OptDial sets dial options for a request, these must be done all at once.
func OptDial(opts ...DialOption) Option {
	return func(r *Request) error {
		if r.Client == nil {
			r.Client = &http.Client{}
		}
		if r.Client.Transport == nil {
			r.Client.Transport = &http.Transport{}
		}
		if typed, ok := r.Client.Transport.(*http.Transport); ok {
			dialer := &net.Dialer{}
			for _, opt := range opts {
				opt(dialer)
			}
			typed.Dial = dialer.Dial
		}
		return nil
	}
}
