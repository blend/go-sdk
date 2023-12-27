/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stringutil

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestSplitSpaceQuoted(t *testing.T) {
	assert := assert.New(t)

	values := SplitSpaceQuoted("")
	assert.Len(values, 0)

	values = SplitSpaceQuoted("foo")
	assert.Len(values, 1)
	assert.Equal("foo", values[0])

	values = SplitSpaceQuoted("foo bar")
	assert.Len(values, 2)
	assert.Equal("foo", values[0])
	assert.Equal("bar", values[1])

	values = SplitSpaceQuoted("foo  bar")
	assert.Len(values, 2)
	assert.Equal("foo", values[0])
	assert.Equal("bar", values[1])

	values = SplitSpaceQuoted("foo\tbar")
	assert.Len(values, 2)
	assert.Equal("foo", values[0])
	assert.Equal("bar", values[1])

	values = SplitSpaceQuoted("foo \tbar")
	assert.Len(values, 2)
	assert.Equal("foo", values[0])
	assert.Equal("bar", values[1])

	values = SplitSpaceQuoted("foo bar  ")
	assert.Len(values, 2)
	assert.Equal("foo", values[0])
	assert.Equal("bar", values[1])

	values = SplitSpaceQuoted("foo bar baz")
	assert.Len(values, 3)
	assert.Equal("foo", values[0])
	assert.Equal("bar", values[1])
	assert.Equal("baz", values[2])

	values = SplitSpaceQuoted(`foo "bar baz"`)
	assert.Len(values, 2, fmt.Sprintf("%#v", values))
	assert.Equal("foo", values[0])
	assert.Equal(`bar baz`, values[1])

	values = SplitSpaceQuoted(`foo --config="bar baz"`)
	assert.Len(values, 2)
	assert.Equal("foo", values[0])
	assert.Equal(`--config="bar baz"`, values[1])

	values = SplitSpaceQuoted(`foo --config='bar baz="hi"'`)
	assert.Len(values, 2)
	assert.Equal("foo", values[0])
	assert.Equal(`--config='bar baz="hi"'`, values[1])

	values = SplitSpaceQuoted(`foo --config="bar baz='hi'"`)
	assert.Len(values, 2)
	assert.Equal("foo", values[0])
	assert.Equal(`--config="bar baz='hi'"`, values[1])

	values = SplitSpaceQuoted(`“what is xyz’s best kept secret?” “1st floor bathrooms” “xyz street side exit” “other (please explain)”`)
	assert.Len(values, 4)
	assert.Equal(`what is xyz’s best kept secret?`, values[0])
	assert.Equal(`1st floor bathrooms`, values[1])
	assert.Equal(`xyz street side exit`, values[2])
	assert.Equal(`other (please explain)`, values[3])

	values = SplitSpaceQuoted(`‘how about now’ ‘or now’`)
	assert.Len(values, 2)
	assert.Equal(`how about now`, values[0])
	assert.Equal(`or now`, values[1])
}
