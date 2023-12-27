/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package mathutil

// Min finds the lowest value in a slice.
func Min(input []float64) float64 {
	if len(input) == 0 {
		return 0
	}

	min := input[0]
	for i := 1; i < len(input); i++ {
		if input[i] < min {
			min = input[i]
		}
	}
	return min
}

// MinInts finds the lowest value in a slice.
func MinInts(input []int) int {
	if len(input) == 0 {
		return 0
	}

	min := input[0]
	for i := 1; i < len(input); i++ {
		if input[i] < min {
			min = input[i]
		}
	}
	return min
}
