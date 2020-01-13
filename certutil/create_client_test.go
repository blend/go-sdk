package certutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
)

func TestCreateClient(t *testing.T) {
	assert := assert.New(t)

	ca, err := CreateCertificateAuthority()
	assert.Nil(err)

	uid := uuid.V4().String()
	client, err := CreateClient(uid, ca)
	assert.Nil(err)
	assert.Len(client.Certificates, 2)
	assert.Len(client.CertificateDERs, 2)
}
