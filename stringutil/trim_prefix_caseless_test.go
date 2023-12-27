/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
