/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package mathutil

import (
	"time"
)

// Sum adds all the numbers of a slice together
func Sum(input []float64) float64 {
	var total float64
	if len(input) == 0 {
		return 0
	}
	// Add em up
	for _, n := range input {
		total += n
	}

	return total
}

// SumInts adds all the numbers of a slice together
func SumInts(values []int) int {
	var total int
	for x := 0; x < len(values); x++ {
		total += values[x]
	}

	return total
}

// SumDurations adds all the numbers of a slice together
func SumDurations(values []time.Duration) time.Duration {
	var total time.Duration
	for x := 0; x < len(values); x++ {
		total += values[x]
	}

	return total
}
