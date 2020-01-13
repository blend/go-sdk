package certutil

import (
	"crypto/x509"
	"github.com/blend/go-sdk/ex"
)

// ExtendSystemCertPool extends the system ca pool with a given list of ca cert key pairs.
func ExtendSystemCertPool(keyPairs ...KeyPair) (*x509.CertPool, error) {
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
