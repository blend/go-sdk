package r2

import (
	"crypto/tls"
)

// OptTLSClientCert adds a client cert and key to the request.
func OptTLSClientCert(cert, key []byte) Option {
	return func(r *Request) error {
		transport, err := EnsureHTTPTransport(r)
		if err != nil {
			return err
		}

		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{}
		}
		cert, err := tls.X509KeyPair(cert, key)
		if err != nil {
			return err
		}
		transport.TLSClientConfig.Certificates = append(transport.TLSClientConfig.Certificates, cert)
		return nil
	}
}
