package r2

import (
	"crypto/tls"
	"crypto/x509"
)

// OptTLSRootCAs sets the client tls root ca pool.
func OptTLSRootCAs(pool *x509.CertPool) Option {
	return func(r *Request) error {
		transport, err := EnsureHTTPTransport(r)
		if err != nil {
			return err
		}
		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{}
		}
		transport.TLSClientConfig.RootCAs = pool
		return nil
	}
}
