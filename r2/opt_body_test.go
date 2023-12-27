/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package r2

import (
	"bytes"
	"io"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptBody(t *testing.T) {
	assert := assert.New(t)

	req := New("https://foo.bar.local")

	assert.Nil(OptBody(io.NopCloser(bytes.NewBufferString("this is only a test")))(req))
	assert.NotNil(req.Request.Body)

	contents, err := io.ReadAll(req.Request.Body)
	assert.Nil(err)
	assert.Equal("this is only a test", string(contents))
}
