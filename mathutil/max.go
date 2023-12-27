/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package mathutil

// Max finds the highest value in a slice.
func Max(input []float64) float64 {
	if len(input) == 0 {
		return 0
	}

	max := input[0]

	for i := 1; i < len(input); i++ {
		if input[i] > max {
			max = input[i]
		}
	}

	return max
}

// MaxInts finds the highest value in a slice.
func MaxInts(input []int) int {
	if len(input) == 0 {
		return 0
	}

	max := input[0]

	for i := 1; i < len(input); i++ {
		if input[i] > max {
			max = input[i]
		}
	}

	return max
}
