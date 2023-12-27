/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stringutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestSplitCSV(t *testing.T) {
	assert := assert.New(t)

	assert.Empty(SplitCSV(""))
	assert.Equal([]string{"foo"}, SplitCSV("foo"))
	assert.Equal([]string{"foo", "bar"}, SplitCSV("foo,bar"))
	assert.Equal([]string{"foo", " bar"}, SplitCSV("foo, bar"))
	assert.Equal([]string{" foo ", " bar "}, SplitCSV(" foo , bar "))
	assert.Equal([]string{"foo", "bar", "baz"}, SplitCSV("foo,bar,baz"))
	assert.Equal([]string{"foo", " bar"}, SplitCSV(`foo," bar"`))
	assert.Equal([]string{" foo ", " bar "}, SplitCSV(`" foo "," bar "`))
	assert.Equal([]string{"foo", "bar", "baz,buzz"}, SplitCSV(`foo,bar,"baz,buzz"`))
	assert.Equal([]string{"foo", "bar", "baz,'buzz'"}, SplitCSV(`foo,bar,"baz,'buzz'"`))
	assert.Equal([]string{"foo", "bar", "baz,'buzz"}, SplitCSV(`foo,bar,"baz,'buzz"`))
	// A double quote can be escaped in CSV by doubling it, e.g. `a"b` -> `"a""b"`
	assert.Equal([]string{"foo", "bar", `baz,"buzz"`}, SplitCSV(`foo,bar,"baz,""buzz"""`))
	assert.Equal([]string{"foo", "bar", `baz"buzz"`}, SplitCSV(`foo,bar,"baz""buzz"""`))
}
