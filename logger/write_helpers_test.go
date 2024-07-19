/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package logger

import (
	"testing"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/assert"
)

func TestFormatLabels(t *testing.T) {
	assert := assert.New(t)

	tf := NewTextOutputFormatter(OptTextNoColor())
	actual := FormatLabels(tf, ansi.ColorBlue, Labels{"foo": "bar", "moo": "loo"})
	assert.Equal("foo=bar moo=loo", actual)

	actual = FormatLabels(tf, ansi.ColorBlue, Labels{"moo": "loo", "foo": "bar"})
	assert.Equal("foo=bar moo=loo", actual)

	tf = NewTextOutputFormatter()
	actual = FormatLabels(tf, ansi.ColorBlue, Labels{"foo": "bar", "moo": "loo"})
	assert.Equal(ansi.ColorBlue.Apply("foo")+"=bar "+ansi.ColorBlue.Apply("moo")+"=loo", actual)

	actual = FormatLabels(tf, ansi.ColorBlue, Labels{"moo": "loo", "foo": "bar"})
	assert.Equal(ansi.ColorBlue.Apply("foo")+"=bar "+ansi.ColorBlue.Apply("moo")+"=loo", actual)
}
