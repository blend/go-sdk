/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
