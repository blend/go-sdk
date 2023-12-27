/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stringutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestEqualsCaseless(t *testing.T) {
	assert := assert.New(t)
	assert.True(EqualsCaseless("foo", "FOO"))
	assert.True(EqualsCaseless("foo123", "FOO123"))
	assert.True(EqualsCaseless("!foo123", "!foo123"))
	assert.False(EqualsCaseless("foo", "bar"))
}
