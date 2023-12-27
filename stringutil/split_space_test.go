/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stringutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestSplitSpace(t *testing.T) {
	assert := assert.New(t)

	values := SplitSpace("")
	assert.Len(values, 0)

	values = SplitSpace("foo")
	assert.Len(values, 1)
	assert.Equal("foo", values[0])

	values = SplitSpace("foo bar")
	assert.Len(values, 2)
	assert.Equal("foo", values[0])
	assert.Equal("bar", values[1])

	values = SplitSpace("foo  bar")
	assert.Len(values, 2)
	assert.Equal("foo", values[0])
	assert.Equal("bar", values[1])

	values = SplitSpace("foo\tbar")
	assert.Len(values, 2)
	assert.Equal("foo", values[0])
	assert.Equal("bar", values[1])

	values = SplitSpace("foo \tbar")
	assert.Len(values, 2)
	assert.Equal("foo", values[0])
	assert.Equal("bar", values[1])

	values = SplitSpace("foo bar  ")
	assert.Len(values, 2)
	assert.Equal("foo", values[0])
	assert.Equal("bar", values[1])

	values = SplitSpace("foo bar baz")
	assert.Len(values, 3)
	assert.Equal("foo", values[0])
	assert.Equal("bar", values[1])
	assert.Equal("baz", values[2])
}
