package certutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCreateServer(t *testing.T) {
	assert := assert.New(t)

	authority, err := CreateCertificateAuthority()
	assert.Nil(err)

	assert.Len(authority.CertificateDERs, 1)
	assert.Len(authority.Certificates, 1)

	server, err := CreateServer("warden-server", authority, OptDNSNames("warden-server-test"))
	assert.Nil(err)
	assert.Len(server.Certificates, 2)
	assert.Len(server.CertificateDERs, 2)
	assert.Len(server.Certificates[0].DNSNames, 1)
}
