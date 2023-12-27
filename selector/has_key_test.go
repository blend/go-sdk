/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package selector

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestHasKey(t *testing.T) {
	assert := assert.New(t)

	valid := Labels{
		"foo": "far",
	}
	assert.True(HasKey("foo").Matches(valid))
	assert.False(HasKey("zoo").Matches(valid))
	assert.Equal("foo", HasKey("foo").String())
}
