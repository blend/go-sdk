/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package selector

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNotEquals(t *testing.T) {
	assert := assert.New(t)

	valid := Labels{
		"foo":	"far",
		"moo":	"bar",
	}
	assert.False(NotEquals{Key: "foo", Value: "far"}.Matches(valid))
	assert.True(NotEquals{Key: "zoo", Value: "buzz"}.Matches(valid))
	assert.True(NotEquals{Key: "foo", Value: "bar"}.Matches(valid))
	assert.Equal("foo != bar", NotEquals{Key: "foo", Value: "bar"}.String())
}
