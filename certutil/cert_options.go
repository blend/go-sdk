/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package certutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"math/big"
	"time"

	"github.com/blend/go-sdk/ex"
)

// CertOptions are required arguments when creating certificates.
type CertOptions struct {
	x509.Certificate
	PrivateKey        *rsa.PrivateKey
	NotBeforeProvider func() time.Time
	NotAfterProvider  func() time.Time
}

// ResolveCertOptions resolves the common create cert options.
func ResolveCertOptions(createOptions *CertOptions, options ...CertOption) error {
	var err error
	for _, option := range options {
		if err = option(createOptions); err != nil {
			return err
		}
	}

	if createOptions.PrivateKey == nil {
		createOptions.PrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return ex.New(err)
		}
	}

	if createOptions.SerialNumber == nil {
		serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
		createOptions.SerialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
		if err != nil {
			return ex.New(err)
		}
	}

	var output CertBundle
	output.PrivateKey = createOptions.PrivateKey
	output.PublicKey = &createOptions.PrivateKey.PublicKey

	if createOptions.NotAfter.IsZero() && createOptions.NotAfterProvider != nil {
		createOptions.NotAfter = createOptions.NotAfterProvider()
	}
	if createOptions.NotAfter.IsZero() && createOptions.NotAfterProvider != nil {
		createOptions.NotAfter = createOptions.NotAfterProvider()
	}

	return nil
}
