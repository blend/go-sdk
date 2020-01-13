package secrets

import (
	"crypto/x509"
	"io/ioutil"
	"github.com/blend/go-sdk/ex"
)

// NewCertPool creates a new cert pool.
// This cert pool starts with the system certs.
func NewCertPool() (*CertPool, error) {
	system, err := x509.SystemCertPool()
	if err != nil {
		return nil, ex.New(err)
	}
	return &CertPool{
		pool: system,
	}, nil
}

// CertPool is a wrapper for x509.CertPool.
type CertPool struct {
	pool *x509.CertPool
}

// Pool returns the underlying cert pool.
func (cp *CertPool) Pool() *x509.CertPool {
	return cp.pool
}

// AddPaths adds a ca path to the cert pool.
func (cp *CertPool) AddPaths(paths ...string) error {
	for _, path := range paths {
		cert, err := ioutil.ReadFile(path)
		if err != nil {
			return ex.New(err)
		}
		if ok := cp.pool.AppendCertsFromPEM(cert); !ok {
			return ex.New("append cert failed", ex.OptMessagef("cert path: %s", path))
		}
	}
	return nil
}
