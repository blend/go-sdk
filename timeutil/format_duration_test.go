/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package timeutil

import (
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestFormatDuration(t *testing.T) {
	assert := assert.New(t)

	testCases := [...]struct {
		Input    time.Duration
		Expected string
	}{
		{Input: ((10 * time.Hour) + (9 * time.Minute) + (8 * time.Second) + (7 * time.Millisecond) + (6 * time.Microsecond) + (5 * time.Nanosecond)), Expected: "10h"},
		{Input: ((9 * time.Minute) + (8 * time.Second) + (7 * time.Millisecond) + (6 * time.Microsecond) + (5 * time.Nanosecond)), Expected: "9m"},
		{Input: ((8 * time.Second) + (7 * time.Millisecond) + (6 * time.Microsecond) + (5 * time.Nanosecond)), Expected: "8s"},
		{Input: ((7 * time.Millisecond) + (6 * time.Microsecond) + (5 * time.Nanosecond)), Expected: "7ms"},
		{Input: ((6 * time.Microsecond) + (5 * time.Nanosecond)), Expected: "6µs"},
	}

	for _, tc := range testCases {
		assert.Equal(tc.Expected, FormatDuration(tc.Input), fmt.Sprintf("%#v", tc))
	}
}
