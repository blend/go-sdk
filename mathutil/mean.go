/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package mathutil

import "time"

// Mean gets the average of a slice of numbers
func Mean(input []float64) float64 {
	if len(input) == 0 {
		return 0
	}

	sum := Sum(input)
	return sum / float64(len(input))
}

// MeanInts gets the average of a slice of numbers
func MeanInts(input []int) float64 {
	if len(input) == 0 {
		return 0
	}
	sum := SumInts(input)
	return float64(sum) / float64(len(input))
}

// MeanDurations gets the average of a slice of numbers
func MeanDurations(input []time.Duration) time.Duration {
	if len(input) == 0 {
		return 0
	}

	sum := SumDurations(input)
	mean := uint64(sum) / uint64(len(input))
	return time.Duration(mean)
}
