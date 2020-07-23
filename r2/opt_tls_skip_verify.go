package r2

import (
	"crypto/tls"
)

// OptTLSSkipVerify sets if we should skip verification.
func OptTLSSkipVerify(skipVerify bool) Option {
	return func(r *Request) error {
		transport, err := EnsureHTTPTransport(r)
		if err != nil {
			return err
		}
		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{}
		}
		transport.TLSClientConfig.InsecureSkipVerify = skipVerify
		return nil
	}
}
