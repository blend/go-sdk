package certutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestParseCertPEM(t *testing.T) {
	assert := assert.New(t)

	certs, err := ParseCertPEM(certLiteral)
	assert.Nil(err)
	assert.Len(certs, 2)
	assert.Equal(certLiteralCommonName, certs[0].Subject.CommonName)
}

func TestCommonNamesForCertPair(t *testing.T) {
	assert := assert.New(t)

	kp := KeyPair{Cert: string(certLiteral)}

	certPEM, err := kp.CertBytes()
	assert.Nil(err)

	commonNames, err := CommonNamesForCertPEM(certPEM)
	assert.Nil(err)
	assert.NotEmpty(commonNames)
	assert.Equal(certLiteralCommonName, commonNames[0])
}
