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

func TestTrimSuffixCaseless(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("abc", TrimSuffixCaseless("abcdef", "def"))
	assert.Equal("ab2", TrimSuffixCaseless("ab2def", "DEF"))
	assert.Equal("ab3", TrimSuffixCaseless("ab3DEF", "def"))
	assert.Equal("abcdef", TrimSuffixCaseless("abcdef", "foo"))
	assert.Equal("abc", TrimSuffixCaseless("abc", "abcdef"))
}
