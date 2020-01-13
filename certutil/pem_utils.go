package certutil

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"

	"github.com/blend/go-sdk/ex"
)

// CommonNamesForCertPEM returns the common names from a cert pair.
func CommonNamesForCertPEM(certPEM []byte) ([]string, error) {
	certs, err := ParseCertPEM(certPEM)
	if err != nil {
		return nil, err
	}
	output := make([]string, len(certs))
	for index, cert := range certs {
		output[index] = cert.Subject.CommonName
	}
	return output, nil
}

// ParseCertPEM parses the cert portion of a cert pair.
func ParseCertPEM(certPem []byte) (output []*x509.Certificate, err error) {
	for len(certPem) > 0 {
		var block *pem.Block
		block, certPem = pem.Decode(certPem)
		if block == nil {
			break
		}
		if block.Type != BlockTypeCertificate || len(block.Headers) != 0 {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			err = ex.New(err)
			continue
		}
		output = append(output, cert)
	}

	return
}

// ReadPrivateKeyPEMFromPath reads a private key pem from a given path.
func ReadPrivateKeyPEMFromPath(keyPath string) (*rsa.PrivateKey, error) {
	contents, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, ex.New(err, ex.OptMessagef("key path: %s", keyPath))
	}
	data, _ := pem.Decode(contents)
	pk, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		return nil, ex.New(err)
	}
	return pk, nil
}
