/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stringutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestTrimPrefixCaseless(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("def", TrimPrefixCaseless("abcdef", "abc"))
	assert.Equal("def", TrimPrefixCaseless("abcdef", "ABC"))
	assert.Equal("DEF", TrimPrefixCaseless("abcDEF", "abc"))
	assert.Equal("abcdef", TrimPrefixCaseless("abcdef", "foo"))
	assert.Equal("abc", TrimPrefixCaseless("abc", "abcdef"))
}
