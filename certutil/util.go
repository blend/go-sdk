package certutil

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"

	"github.com/blend/go-sdk/ex"
)

// Errors
const (
	ErrInvalidCertPEM ex.Class = "failed to add cert to pool as pem"
)

// MustBytes panics on an error or returns the contents.
func MustBytes(contents []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return contents
}

// BytesWithError returns a bytes error response with the error
// as an ex.
func BytesWithError(bytes []byte, err error) ([]byte, error) {
	return bytes, ex.New(err)
}

// ReadFiles reads a list of files as bytes.
func ReadFiles(files ...string) (data [][]byte, err error) {
	var contents []byte
	for _, path := range files {
		contents, err = ioutil.ReadFile(path)
		if err != nil {
			return nil, ex.New(err)
		}
		data = append(data, contents)
	}
	return
}

// ExtendSystemPoolWithKeyPairCerts extends the system ca pool with a given list of ca cert key pairs.
func ExtendSystemPoolWithKeyPairCerts(keyPairs ...KeyPair) (*x509.CertPool, error) {
	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, ex.New(err)
	}
	var contents []byte
	for _, keyPair := range keyPairs {
		contents, err = keyPair.CertBytes()
		if err != nil {
			return nil, ex.New(err)
		}
		if ok := pool.AppendCertsFromPEM(contents); !ok {
			return nil, ex.New(ErrInvalidCertPEM)
		}
	}

	return pool, nil
}

// ExtendEmptyPoolWithKeyPairCerts extends an empty pool with a given set of certs.
func ExtendEmptyPoolWithKeyPairCerts(keyPairs ...KeyPair) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	var err error
	var contents []byte
	for _, keyPair := range keyPairs {
		contents, err = keyPair.CertBytes()
		if err != nil {
			return nil, err
		}
		if ok := pool.AppendCertsFromPEM(contents); !ok {
			return nil, ex.New(ErrInvalidCertPEM)
		}
	}
	return pool, nil
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
