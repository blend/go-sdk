/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestOptHost(t *testing.T) {
	assert := assert.New(t)

	r := New(TestURL, OptHost("bar.invalid"))

	assert.Equal("https://bar.invalid/test?query=value", r.Request.URL.String())
}
