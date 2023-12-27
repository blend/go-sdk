/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package bufferutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestPutOnClose(t *testing.T) {
	assert := assert.New(t)

	pool := NewPool(32)

	poc := PutOnClose(pool.Get(), pool)
	assert.NotNil(poc)
	assert.NotNil(poc.Buffer)
	assert.NotNil(poc.Pool)
	assert.Nil(poc.Close())
}
