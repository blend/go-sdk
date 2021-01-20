package certutil

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/blend/go-sdk/ex"
)

// NewClientTLSConfig returns a new client tls config.
// This is useful for making mutual tls calls to servers that require it.
func NewClientTLSConfig(clientCert KeyPair, certificateAuthorities []KeyPair) (*tls.Config, error) {
	clientCertPEM, err := clientCert.CertBytes()
	if err != nil {
		return nil, ex.New(err)
	}
	clientKeyPEM, err := clientCert.KeyBytes()
	if err != nil {
		return nil, ex.New(err)
	}

	if len(clientCertPEM) == 0 {
		return nil, ex.New("invalid key pair; empty cert pem data")
	}
	if len(clientKeyPEM) == 0 {
		return nil, ex.New("invalid key pair; empty key pem data")
	}
	cert, err := tls.X509KeyPair(clientCertPEM, clientKeyPEM)
	if err != nil {
		return nil, ex.New(err)
	}
	config := new(tls.Config)

	rootCAPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, ex.New(err)
	}

	for _, caCert := range certificateAuthorities {
		contents, err := caCert.CertBytes()
		if err != nil {
			return nil, ex.New(err)
		}
		if ok := rootCAPool.AppendCertsFromPEM(contents); !ok {
			return nil, ex.New("failed to append ca cert file")
		}
	}

	config.Certificates = []tls.Certificate{cert}
	config.RootCAs = rootCAPool
	return config, nil
}
