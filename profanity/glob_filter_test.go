/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package profanity

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_Glob_include(t *testing.T) {
	its := assert.New(t)

	testCases := [...]struct {
		Input    string
		Expected bool
	}{
		{Input: "foo.txt", Expected: false},
		{Input: "foo.go", Expected: true},
		{Input: "Dockerfile", Expected: true},
		{Input: "Dockerfile.unit-test", Expected: true},
		{Input: "project/path/foo.txt", Expected: false},
		{Input: "project/path/foo.go", Expected: true},
		{Input: "project/path/Dockerfile", Expected: true},
		{Input: "project/path/testing.Dockerfile", Expected: false},
		{Input: "project/path/Dockerfile.unit-test", Expected: true},
		{Input: "project/path/Dockerfiles/testing-file", Expected: false},
	}

	filter := GlobFilter{
		Filter: Filter{
			Include: []string{"*.go", "Dockerfile", "Dockerfile.*", "**/Dockerfile", "**/Dockerfile.*"},
		},
	}

	var actual bool
	for _, tc := range testCases {
		actual = filter.Allow(tc.Input)
		its.Equal(tc.Expected, actual, fmt.Sprintf("%s should yield %v", tc.Input, tc.Expected))
	}
}
