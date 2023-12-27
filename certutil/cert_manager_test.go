/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package certutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
)

func TestNewCertManagerWithKeyPairs(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	mgr, err := NewCertManagerWithKeyPairs(KeyPair{
		CertPath:	"testdata/server.cert.pem",
		KeyPath:	"testdata/server.key.pem",
	}, []KeyPair{{
		CertPath: "testdata/ca.cert.pem",
	}}, KeyPair{
		CertPath: "testdata/client.cert.pem",
	})
	assert.Nil(err)
	assert.NotEmpty(mgr.ClientCerts)
	assert.NotNil(mgr.TLSConfig)
	assert.NotNil(mgr.TLSConfig.RootCAs)
}

func TestCertManagerAddClientCert(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	mgr := NewCertManager()
	assert.Nil(mgr.AddClientCert(certLiteral))
	assert.NotEmpty(mgr.ClientCerts)
	assert.True(mgr.HasClientCert(certLiteralCommonName))
	assert.NotNil(mgr.AddClientCert(uuid.V4()))
}

func TestCertManagerRemoveClientCert(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	mgr := NewCertManager()
	assert.Nil(mgr.AddClientCert(certLiteral))
	assert.NotEmpty(mgr.ClientCerts)

	assert.Nil(mgr.RemoveClientCert(certLiteralCommonName))
	assert.Empty(mgr.ClientCerts)
}

func TestCertManagerGetConfigForClient(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	mgr := NewCertManager()
	assert.Nil(mgr.AddClientCert(certLiteral))
	assert.NotEmpty(mgr.ClientCerts)

	config, err := mgr.GetConfigForClient(nil)
	assert.Nil(err)
	assert.NotEmpty(config.ClientCAs.Subjects())
}
