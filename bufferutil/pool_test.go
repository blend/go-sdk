/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package bufferutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestPool(t *testing.T) {
	assert := assert.New(t)

	pool := NewPool(1024)
	buf := pool.Get()
	assert.NotNil(buf)
	assert.Equal(1024, buf.Cap())
	assert.Zero(buf.Len())
	pool.Put(buf)
}
