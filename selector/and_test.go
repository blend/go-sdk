/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package selector

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestAnd(t *testing.T) {
	assert := assert.New(t)

	valid := Labels{
		"foo": "far",
		"moo": "lar",
	}
	invalid := Labels{
		"foo": "far",
		"moo": "bar",
	}

	selector := And([]Selector{Equals{Key: "foo", Value: "far"}, Equals{Key: "moo", Value: "lar"}})
	assert.True(selector.Matches(valid))
	assert.False(selector.Matches(invalid))

	assert.Equal("foo == far, moo == lar", selector.String())
}
