/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package certutil

import (
	"crypto/x509"

	"github.com/blend/go-sdk/ex"
)

// CreateCertPool extends an empty pool with a given set of certs.
func CreateCertPool(keyPairs ...KeyPair) (*x509.CertPool, error) {
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
