/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package mathutil

// Median gets the median number in a slice of numbers
func Median(input []float64) (median float64) {
	l := len(input)
	if l == 0 {
		return 0
	}

	median = MedianSorted(CopySort(input))
	return
}

// MedianSorted gets the median number in a sorted slice of numbers
func MedianSorted(sortedInput []float64) (median float64) {
	l := len(sortedInput)
	if l == 0 {
		return 0
	}

	if l%2 == 0 {
		median = (sortedInput[(l>>1)-1] + sortedInput[l>>1]) / 2.0
	} else {
		median = sortedInput[l>>1]
	}

	return median
}
