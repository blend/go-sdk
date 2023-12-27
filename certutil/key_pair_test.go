/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package certutil

import (
	"os"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestKeyPairIsZero(t *testing.T) {
	assert := assert.New(t)

	assert.True(KeyPair{}.IsZero())
	assert.False(KeyPair{Cert: "foo"}.IsZero())
	assert.False(KeyPair{Key: "foo"}.IsZero())
	assert.False(KeyPair{CertPath: "foo"}.IsZero())
	assert.False(KeyPair{KeyPath: "foo"}.IsZero())
}

func TestKeyPairCertBytes(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("foo", MustBytes(KeyPair{Cert: "foo"}.CertBytes()))
	assert.Equal(MustBytes(os.ReadFile("testdata/client.cert.pem")), MustBytes(KeyPair{CertPath: "testdata/client.cert.pem"}.CertBytes()))
}

func TestKeyPairKeyBytes(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("foo", MustBytes(KeyPair{Key: "foo"}.KeyBytes()))
	assert.Equal(MustBytes(os.ReadFile("testdata/client.key.pem")), MustBytes(KeyPair{KeyPath: "testdata/client.key.pem"}.KeyBytes()))
}

func TestKeyPairString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("[ cert: <literal>, key: <literal> ]", KeyPair{Cert: "bar", Key: "foo"}.String())
	assert.Equal("[ cert: bar, key: foo ]", KeyPair{CertPath: "bar", KeyPath: "foo"}.String())
}

func TestKeyPairTLSCertificate(t *testing.T) {
	its := assert.New(t)

	kp := KeyPair{
		CertPath:	"testdata/server.cert.pem",
		KeyPath:	"testdata/server.key.pem",
	}

	cert, err := kp.TLSCertificate()
	its.Nil(err)
	its.NotNil(cert)
}
