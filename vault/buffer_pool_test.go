/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package vault

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestBufferPool(t *testing.T) {
	assert := assert.New(t)

	pool := NewBufferPool(1024)
	buf := pool.Get()
	assert.NotNil(buf)
	pool.Put(buf)
}
