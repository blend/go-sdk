/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package certutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestExtendSystemCertPool(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	pool, err := ExtendSystemCertPool(KeyPair{Cert: string(caCertLiteral)})
	assert.Nil(err)
	assert.NotNil(pool)
}
