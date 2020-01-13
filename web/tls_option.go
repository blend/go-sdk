package web

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/blend/go-sdk/ex"
)

// TLSOption is an option for TLS configs.
type TLSOption func(*tls.Config) error

// OptTLSClientCertPool adds a given set of certs in binary PEM format
// to the system CA pool.
func OptTLSClientCertPool(certPEMs ...[]byte) TLSOption {
	return func(t *tls.Config) error {
		if t == nil {
			t = &tls.Config{}
		}
		t.ClientCAs = x509.NewCertPool()
		for _, certPem := range certPEMs {
			ok := t.ClientCAs.AppendCertsFromPEM(certPem)
			if !ok {
				return ex.New("invalid ca cert for client cert pool")
			}
		}
		t.BuildNameToCertificate()

		// this forces the server to reload the tls config for every request if there is a cert pool loaded.
		// normally this would introduce overhead but it allows us to hot patch the cert pool.
		t.GetConfigForClient = func(_ *tls.ClientHelloInfo) (*tls.Config, error) {
			return t, nil
		}
		return nil
	}
}

// OptTLSClientCertVerification sets the verification level for client certs.
func OptTLSClientCertVerification(verification tls.ClientAuthType) TLSOption {
	return func(t *tls.Config) error {
		if t == nil {
			t = &tls.Config{}
		}
		t.ClientAuth = verification
		return nil
	}
}
