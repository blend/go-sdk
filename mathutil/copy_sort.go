/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package mathutil

import (
	"sort"
	"time"
)

// CopySort copies and sorts an array of floats.
func CopySort(input []float64) []float64 {
	copy := Copy(input)
	sort.Float64s(copy)
	return copy
}

// Copy copies an array of float64s.
func Copy(input []float64) []float64 {
	output := make([]float64, len(input))
	copy(output, input)
	return output
}

// CopySortInts copies and sorts an array of floats.
func CopySortInts(input []int) []int {
	copy := CopyInts(input)
	sort.Ints(copy)
	return copy
}

// CopyInts copies an array of float64s.
func CopyInts(input []int) []int {
	output := make([]int, len(input))
	copy(output, input)
	return output
}

// CopySortDurations copies and sorts an array of floats.
func CopySortDurations(input []time.Duration) []time.Duration {
	copy := CopyDurations(input)
	sort.Sort(Durations(copy))
	return copy
}

// CopyDurations copies an array of time.Duration.
func CopyDurations(input []time.Duration) []time.Duration {
	output := make([]time.Duration, len(input))
	copy(output, input)
	return output
}

// Durations is an array of durations.
type Durations []time.Duration

// Len implements sort.Sorter
func (d Durations) Len() int {
	return len(d)
}

// Swap implements sort.Sorter
func (d Durations) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Less implements sort.Sorter
func (d Durations) Less(i, j int) bool {
	return d[i] < d[j]
}
