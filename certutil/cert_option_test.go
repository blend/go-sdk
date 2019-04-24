package certutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptSubjectName(t *testing.T) {
	assert := assert.New(t)
	// create the ca
	ca, err := CreateCertificateAuthority()
	assert.Nil(err)

	// create the server certs
	server, err := CreateServer("mtls-example.local", ca, OptSubjectCommonName("localhost"))
	assert.Nil(err)
	names, err := server.CommonNames()
	assert.Nil(err)
	assert.Equal([]string{"localhost"}, names)
}
