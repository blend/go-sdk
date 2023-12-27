/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package grpcutil

import (
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

// scale duration by a factor
func scaleDuration(d time.Duration, factor float64) time.Duration {
	return time.Duration(float64(d) * factor)
}

func TestJitterUp(t *testing.T) {
	assert := assert.New(t)

	// arguments to jitterup
	duration := 10 * time.Second
	variance := 0.10

	// bound to check
	max := 11000 * time.Millisecond
	min := 9000 * time.Millisecond
	high := scaleDuration(max, 0.98)
	low := scaleDuration(min, 1.02)

	highCount := 0
	lowCount := 0

	for i := 0; i < 1000; i++ {
		out := JitterUp(duration, variance)
		assert.True(out <= max, fmt.Sprintf("value %s must be <= %s", out, max))
		assert.True(out >= min, fmt.Sprintf("value %s must be >= %s", out, min))

		if out > high {
			highCount++
		}
		if out < low {
			lowCount++
		}
	}

	assert.True(highCount != 0, fmt.Sprintf("at least one sample should reach to > %s", high))
	assert.True(lowCount != 0, fmt.Sprintf("at least one sample should to < %s", low))
}
