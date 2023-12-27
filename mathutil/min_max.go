/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package mathutil

import "time"

// MinMax returns both the min and max in one pass.
func MinMax(values []float64) (min, max float64) {
	if len(values) == 0 {
		return
	}
	min = values[0]
	max = values[0]
	for _, v := range values {
		if max < v {
			max = v
		}
		if min > v {
			min = v
		}
	}
	return
}

// MinMaxInts returns both the min and max of ints in one pass.
func MinMaxInts(values []int) (min, max int) {
	if len(values) == 0 {
		return
	}
	min = values[0]
	max = values[0]
	for _, v := range values {
		if max < v {
			max = v
		}
		if min > v {
			min = v
		}
	}
	return
}

// MinMaxDurations returns both the min and max of time.Duration in one pass.
func MinMaxDurations(values []time.Duration) (min, max time.Duration) {
	if len(values) == 0 {
		return
	}
	min = values[0]
	max = values[0]
	for _, v := range values {
		if max < v {
			max = v
		}
		if min > v {
			min = v
		}
	}
	return
}
