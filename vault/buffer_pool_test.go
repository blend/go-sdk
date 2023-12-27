/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
