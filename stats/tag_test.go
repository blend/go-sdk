/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stats

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_Tag(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		Key      string
		Value    string
		Expected string
	}{
		{Key: "foo", Value: "bar", Expected: "foo:bar"},
		{Key: "foo1", Value: "bar:", Expected: "foo1:bar:"},
		{Key: "foo_", Value: "bar_", Expected: "foo_:bar_"},
		{Key: "foo%", Value: "bar$", Expected: "foo_:bar_"},
		{Key: "foo;", Value: "bar#", Expected: "foo_:bar_"},
		{Key: "foo;", Value: "bar#", Expected: "foo_:bar_"},
	}

	for _, tc := range testCases {
		assert.Equal(tc.Expected, Tag(tc.Key, tc.Value))
	}
}
