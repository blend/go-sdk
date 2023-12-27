/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package r2

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptMethods(t *testing.T) {
	assert := assert.New(t)

	req := New("https://foo.bar.local")

	assert.Nil(OptMethod("OPTIONS")(req))
	assert.Equal("OPTIONS", req.Request.Method)

	assert.Nil(OptGet()(req))
	assert.Equal("GET", req.Request.Method)

	assert.Nil(OptPost()(req))
	assert.Equal("POST", req.Request.Method)

	assert.Nil(OptPut()(req))
	assert.Equal("PUT", req.Request.Method)

	assert.Nil(OptPatch()(req))
	assert.Equal("PATCH", req.Request.Method)

	assert.Nil(OptDelete()(req))
	assert.Equal("DELETE", req.Request.Method)
}
