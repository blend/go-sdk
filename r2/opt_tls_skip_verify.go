package r2

import (
	"crypto/tls"
	"net/http"
)

// OptTLSSkipVerify sets if we should skip verification.
func OptTLSSkipVerify(skipVerify bool) Option {
	return func(r *Request) error {
		if r.Client == nil {
			r.Client = &http.Client{}
		}

		if r.Client.Transport == nil {
			r.Client.Transport = &http.Transport{}
		}
		if typed, ok := r.Client.Transport.(*http.Transport); ok {
			if typed.TLSClientConfig == nil {
				typed.TLSClientConfig = &tls.Config{}
			}
			typed.TLSClientConfig.InsecureSkipVerify = skipVerify
		}

		return nil
	}
}
