/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptUserAgent(t *testing.T) {
	assert := assert.New(t)

	opt := OptUserAgent("blend test harness")
	req := New(TestURL)
	assert.NotEqual("blend test harness", req.Request.UserAgent())
	assert.Nil(opt(req))
	assert.Equal("blend test harness", req.Request.UserAgent())
}
