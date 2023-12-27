/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package certutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCreateCertPool(t *testing.T) {
	assert := assert.New(t)

	pool, err := CreateCertPool(KeyPair{Cert: string(caCertLiteral)})
	assert.Nil(err)
	assert.NotNil(pool)
}
