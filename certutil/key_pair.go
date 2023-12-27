/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package certutil

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/blend/go-sdk/ex"
)

// NewKeyPairFromPaths returns a key pair from paths.
func NewKeyPairFromPaths(certPath, keyPath string) KeyPair {
	return KeyPair{CertPath: certPath, KeyPath: keyPath}
}

// KeyPair is an x509 pem key pair as strings.
type KeyPair struct {
	Cert		string	`json:"cert,omitempty" yaml:"cert,omitempty"`
	CertPath	string	`json:"certPath,omitempty" yaml:"certPath,omitempty"`
	Key		string	`json:"key,omitempty" yaml:"key,omitempty"`
	KeyPath		string	`json:"keyPath,omitempty" yaml:"keyPath,omitempty"`
}

// IsZero returns if the key pair is set or not.
func (kp KeyPair) IsZero() bool {
	return kp.Cert == "" &&
		kp.Key == "" &&
		kp.CertPath == "" &&
		kp.KeyPath == ""
}

// IsCertPath returns if the keypair cert is a path.
func (kp KeyPair) IsCertPath() bool {
	return kp.Cert == "" && kp.CertPath != ""
}

// IsKeyPath returns if the keypair key is a path.
func (kp KeyPair) IsKeyPath() bool {
	return kp.Key == "" && kp.KeyPath != ""
}

// CertBytes returns the key pair cert bytes.
func (kp KeyPair) CertBytes() ([]byte, error) {
	if kp.Cert != "" {
		return []byte(kp.Cert), nil
	}
	if kp.CertPath == "" {
		return nil, ex.New("error loading cert; cert path unset")
	}
	contents, err := os.ReadFile(os.ExpandEnv(kp.CertPath))
	if err != nil {
		return nil, ex.New("error loading cert from path", ex.OptInner(err), ex.OptMessage(kp.CertPath))
	}
	return contents, nil
}

// KeyBytes returns the key pair key bytes.
func (kp KeyPair) KeyBytes() ([]byte, error) {
	if kp.Key != "" {
		return []byte(kp.Key), nil
	}
	if kp.KeyPath == "" {
		return nil, ex.New("error loading key; key path unset")
	}
	contents, err := os.ReadFile(os.ExpandEnv(kp.KeyPath))
	if err != nil {
		return nil, ex.New("error loading key from path", ex.OptInner(err), ex.OptMessage(kp.KeyPath))
	}
	return contents, nil
}

// String returns a string representation of the key pair.
func (kp KeyPair) String() (output string) {
	output = "[ "
	if kp.Cert != "" {
		output += "cert: <literal>"
	} else if kp.CertPath != "" {
		output += ("cert: " + os.ExpandEnv(kp.CertPath))
	}
	if kp.Key != "" {
		output += ", key: <literal>"
	} else if kp.KeyPath != "" {
		output += (", key: " + os.ExpandEnv(kp.KeyPath))
	}
	output += " ]"
	return output
}

// TLSCertificate returns the KeyPair as a tls.Certificate.
func (kp KeyPair) TLSCertificate() (*tls.Certificate, error) {
	certBytes, err := kp.CertBytes()
	if err != nil {
		return nil, ex.New(err)
	}
	keyBytes, err := kp.KeyBytes()
	if err != nil {
		return nil, ex.New(err)
	}
	cert, err := tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		return nil, ex.New(err)
	}
	return &cert, nil
}

// TLSCertificateWithLeaf returns the KeyPair as a tls.Certificate.
func (kp KeyPair) TLSCertificateWithLeaf() (*tls.Certificate, error) {
	cert, err := kp.TLSCertificate()
	if err != nil {
		return nil, err
	}
	if len(cert.Certificate) == 0 {
		return nil, ex.New("invalid certificate; empty certificate list")
	}
	if cert.Leaf == nil {
		cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
		if err != nil {
			return nil, ex.New(err)
		}
	}
	return cert, nil
}
