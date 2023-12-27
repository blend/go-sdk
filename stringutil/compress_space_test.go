/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stringutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCompressSpace(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("", CompressSpace(""))
	assert.Equal("", CompressSpace(" "))
	assert.Equal("", CompressSpace("\n"))
	assert.Equal("", CompressSpace("\t"))

	assert.Equal("foo", CompressSpace(" foo"))
	assert.Equal("foo", CompressSpace("foo "))
	assert.Equal("foo", CompressSpace("foo\n"))

	assert.Equal("foo bar", CompressSpace("foo bar"))
	assert.Equal("foo bar", CompressSpace("foo\tbar"))
	assert.Equal("foo bar", CompressSpace("foo\nbar"))

	assert.Equal("foo bar", CompressSpace("foo  bar"))
	assert.Equal("foo bar", CompressSpace("foo\t\tbar"))
	assert.Equal("foo bar", CompressSpace("foo\n\nbar"))

	assert.Equal("foo bar baz", CompressSpace("foo  bar   baz"))
	assert.Equal("foo bar baz", CompressSpace("foo\t\t\tbar baz\n"))
	assert.Equal("foo bar baz", CompressSpace("foo\n\n\nbar\tbaz"))
}
