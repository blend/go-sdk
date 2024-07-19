/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package consistenthash

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_StableHash_isStable(t *testing.T) {
	its := assert.New(t)

	testCases := [...]struct {
		Input    string
		Expected uint64
	}{
		{Input: "foo-bar-baz", Expected: 0x3bcce3e4ec07ffbc},
		{Input: "google.com", Expected: 0x1c1766d80c8f9809},
		{Input: "worker-5", Expected: 0xd95dff1c56889f11},
		{Input: "worker-5|0", Expected: 0xffbfaa9d0532a241},
	}

	for _, testCase := range testCases {
		its.Equal(testCase.Expected, StableHash([]byte(testCase.Input)))
	}
}
