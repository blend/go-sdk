/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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
