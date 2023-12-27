/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package mathutil

// Mode gets the mode of a slice of numbers
// `Mode` generally is the most frequently occurring values within the input set.
func Mode(input []float64) []float64 {
	l := len(input)
	if l == 1 {
		return input
	} else if l == 0 {
		return []float64{}
	}

	m := make(map[float64]int)
	for _, v := range input {
		m[v]++
	}

	mode := []float64{}

	var current int
	for k, v := range m {
		switch {
		case v < current:
		case v > current:
			current = v
			mode = append(mode[:0], k)
		default:
			mode = append(mode, k)
		}
	}

	lm := len(mode)
	if l == lm {
		return []float64{}
	}

	return mode
}
