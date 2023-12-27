/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stringutil

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_Glob(t *testing.T) {
	its := assert.New(t)

	testCases := [...]struct {
		Subj     string
		Pattern  string
		Expected bool
	}{
		{"", "", true},
		{"test", "", false},
		{"", "false", false},
		{"", "*", true},
		{"foo", "*", true},
		{"bar", "*", true},

		{"bar/foo", "bar/*", true},
		{"foo/bar", "bar/*", false},

		{"bar/loo/foo", "bar/*", true},
		{"foo/bar/loo", "bar/*", false},

		{"foo/bar/baz/buzz", "*/bar/*", true},
		{"/foo/bar/baz/buzz", "*/bar/*", true},
		{"foo/bar/baz/buzz", "*/foo/*", false},
		{"foo/bar/baz/buzz", "foo/*", true},
		{"/foo/bar/baz/buzz", "*foo/*", true},
		{"test", "te*", true},
		{"test", "*st", true},
		{"test", "foo", false},
		{"test", "foo*", false},
		{"test", "*foo*", false},
	}

	for _, testCase := range testCases {
		its.Equal(testCase.Expected, Glob(testCase.Subj, testCase.Pattern), fmt.Sprint(testCase.Subj, " => ", testCase.Pattern))
	}
}
