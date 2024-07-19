/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stringutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCSV(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("", CSV(nil))
	assert.Equal("", CSV([]string{}))
	assert.Equal("foo", CSV([]string{"foo"}))
	assert.Equal("foo,bar", CSV([]string{"foo", "bar"}))
	assert.Equal("foo,bar,baz", CSV([]string{"foo", "bar", "baz"}))
	assert.Equal(`foo," bar"`, CSV([]string{"foo", " bar"}))
	assert.Equal(`" foo "," bar "`, CSV([]string{" foo ", " bar "}))
	assert.Equal(`foo,bar,"baz,buzz"`, CSV([]string{"foo", "bar", "baz,buzz"}))
	assert.Equal(`foo,bar,"baz,'buzz'"`, CSV([]string{"foo", "bar", "baz,'buzz'"}))
	assert.Equal(`foo,bar,"baz,'buzz"`, CSV([]string{"foo", "bar", "baz,'buzz"}))
	// A double quote can be escaped in CSV by doubling it, e.g. `a"b` -> `"a""b"`
	assert.Equal(`foo,bar,"baz,""buzz"""`, CSV([]string{"foo", "bar", `baz,"buzz"`}))
	assert.Equal(`foo,bar,"baz""buzz"""`, CSV([]string{"foo", "bar", `baz"buzz"`}))
}
