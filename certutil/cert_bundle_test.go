package certutil

import (
	"bytes"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNewCertBundleFromLiterals(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
