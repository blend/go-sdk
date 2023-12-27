/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
