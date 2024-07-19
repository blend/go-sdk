/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"strings"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptBasicAuth(t *testing.T) {
	assert := assert.New(t)

	opt := OptBasicAuth("foo", "bar")

	req := New("https://foo.bar.local")
	assert.Nil(opt(req))

	assert.NotNil(req.Request.Header)
	assert.NotEmpty(req.Request.Header.Get("Authorization"))
	assert.True(strings.HasPrefix(req.Request.Header.Get("Authorization"), "Basic "))
}
