/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package certutil

import (
	"bytes"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNewCertBundleFromLiterals(t *testing.T) {
	assert := assert.New(t)

	bundle, err := NewCertBundle(KeyPair{
		Cert: string(certLiteral),
		Key:  string(keyLiteral),
	})

	assert.Nil(err)
	assert.NotNil(bundle.PrivateKey)
	assert.NotNil(bundle.PublicKey)
	assert.NotEmpty(bundle.Certificates)
	assert.NotEmpty(bundle.CertificateDERs)
}

func TestNewCertBundleFromFiles(t *testing.T) {
	assert := assert.New(t)

	bundle, err := NewCertBundle(KeyPair{
		CertPath: "testdata/client.cert.pem",
		KeyPath:  "testdata/client.key.pem",
	})
	assert.Nil(err)
	assert.NotNil(bundle.PrivateKey)
	assert.NotNil(bundle.PublicKey)
	assert.NotEmpty(bundle.Certificates)
	assert.NotEmpty(bundle.CertificateDERs)
}

func TestCertBundleWriteCertPem(t *testing.T) {
	assert := assert.New(t)

	bundle, err := NewCertBundle(KeyPair{
		CertPath: "testdata/client.cert.pem",
		KeyPath:  "testdata/client.key.pem",
	})
	assert.Nil(err)

	buffer := new(bytes.Buffer)
	assert.Nil(bundle.WriteCertPem(buffer))
	assert.NotZero(buffer.Len())
}

func TestCertBundleWriteKeyPem(t *testing.T) {
	assert := assert.New(t)

	bundle, err := NewCertBundle(KeyPair{
		CertPath: "testdata/client.cert.pem",
		KeyPath:  "testdata/client.key.pem",
	})
	assert.Nil(err)

	buffer := new(bytes.Buffer)
	assert.Nil(bundle.WriteKeyPem(buffer))
	assert.NotZero(buffer.Len())
}

func TestCertBundleCommonNames(t *testing.T) {
	assert := assert.New(t)

	bundle, err := NewCertBundle(KeyPair{
		Cert: string(certLiteral),
		Key:  string(keyLiteral),
	})
	assert.Nil(err)

	commonNames, err := bundle.CommonNames()
	assert.Nil(err)
	assert.Len(commonNames, 2)
}

func TestCertBundle_ServerConfig(t *testing.T) {
	assert := assert.New(t)

	bundle, err := NewCertBundle(KeyPair{
		Cert: string(certLiteral),
		Key:  string(keyLiteral),
	})
	assert.Nil(err)

	cfg, err := bundle.ServerConfig()
	assert.Nil(err)
	assert.NotEmpty(cfg.Certificates)
	assert.NotNil(cfg.RootCAs)
}
