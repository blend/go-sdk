/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package selector

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNotIn(t *testing.T) {
	assert := assert.New(t)

	valid := Labels{
		"foo": "mar",
		"moo": "lar",
	}
	invalid := Labels{
		"foo": "far",
		"moo": "lar",
	}
	missing := Labels{
		"loo": "mar",
		"moo": "lar",
	}

	selector := NotIn{Key: "foo", Values: []string{"bar", "far"}}
	assert.True(selector.Matches(valid))
	assert.True(selector.Matches(missing))
	assert.False(selector.Matches(invalid))
	assert.Equal("foo notin (bar, far)", selector.String())
}
