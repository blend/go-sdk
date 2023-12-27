/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package certutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptSubjectName(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	caKeyPair := KeyPair{
		Cert:	string(caCertLiteral),
		Key:	string(caKeyLiteral),
	}
	ca, err := NewCertBundle(caKeyPair)
	assert.Nil(err)

	// create the server certs
	server, err := CreateServer("mtls-example.local", ca, OptSubjectCommonName("localhost"))
	assert.Nil(err)
	names, err := server.CommonNames()
	assert.Nil(err)
	assert.Equal([]string{"localhost", "warden-ca"}, names)
}
