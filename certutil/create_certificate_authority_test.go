package certutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCreateCA(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	ca, err := CreateCertificateAuthority()
	assert.Nil(err)
	assert.NotNil(ca.PrivateKey)
	assert.NotNil(ca.PublicKey)
	assert.Len(ca.Certificates, 1)
	assert.Len(ca.CertificateDERs, 1)
}
