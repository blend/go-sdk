/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package selector

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNotHasKey(t *testing.T) {
	assert := assert.New(t)

	valid := Labels{
		"foo": "far",
	}
	assert.False(NotHasKey("foo").Matches(valid))
	assert.True(NotHasKey("zoo").Matches(valid))
	assert.Equal("!foo", NotHasKey("foo").String())
}
