/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

import (
	"crypto/tls"
	"crypto/x509"
	"os"

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
		for _, certPEM := range certPEMs {
			ok := t.ClientCAs.AppendCertsFromPEM(certPEM)
			if !ok {
				return ex.New("invalid ca cert for client cert pool")
			}
		}
		// this is deprecated
		// t.BuildNameToCertificate()

		// this forces the server to reload the tls config for every request if there is a cert pool loaded.
		// normally this would introduce overhead but it allows us to hot patch the cert pool.
		t.GetConfigForClient = func(_ *tls.ClientHelloInfo) (*tls.Config, error) {
			return t, nil
		}
		return nil
	}
}

// OptTLSAppendCACertsFromPaths adds CA certs from file paths
// to the system CA pool.
func OptTLSAppendCACertsFromPaths(paths ...string) TLSOption {
	return func(t *tls.Config) error {
		if t == nil {
			return nil
		}
		t.ClientCAs = x509.NewCertPool()
		for _, path := range paths {
			ca, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if ok := t.ClientCAs.AppendCertsFromPEM(ca); !ok {
				return ex.New("invalid ca cert for client cert pool")
			}
		}

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
