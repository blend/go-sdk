package certutil

import (
	"io/ioutil"
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
	assert.Equal(MustBytes(ioutil.ReadFile("testdata/client.cert.pem")), MustBytes(KeyPair{CertPath: "testdata/client.cert.pem"}.CertBytes()))
}

func TestKeyPairKeyBytes(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("foo", MustBytes(KeyPair{Key: "foo"}.KeyBytes()))
	assert.Equal(MustBytes(ioutil.ReadFile("testdata/client.key.pem")), MustBytes(KeyPair{KeyPath: "testdata/client.key.pem"}.KeyBytes()))
}

func TestKeyPairString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("[ cert: <literal>, key: <literal> ]", KeyPair{Cert: "bar", Key: "foo"}.String())
	assert.Equal("[ cert: bar, key: foo ]", KeyPair{CertPath: "bar", KeyPath: "foo"}.String())
}
