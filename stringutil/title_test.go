/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stringutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestTitle(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("123456", Title("123456"))
	assert.Equal("Test", Title("test"))
	assert.Equal("Test", Title("TEST"))
	assert.Equal("Test", Title("Test"))
	assert.Equal("Test Strings", Title("test strings"))
	assert.Equal("Test_Strings", Title("test_strings"))
	assert.Equal("Test_Strings", Title("TEST_STRINGS"))
}
