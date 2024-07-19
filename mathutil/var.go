/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package mathutil

// Var finds the variance for both population and sample data
func Var(input []float64, sample int) (variance float64) {
	if len(input) == 0 {
		return 0
	}
	m := Mean(input)

	for _, n := range input {
		variance += (float64(n) - m) * (float64(n) - m)
	}

	// When getting the mean of the squared differences
	// "sample" will allow us to know if it's a sample
	// or population and wether to subtract by one or not
	variance = variance / float64((len(input) - (1 * sample)))
	return
}

// VarP finds the amount of variance within a population
func VarP(input []float64) float64 {
	return Var(input, 0)
}

// VarS finds the amount of variance within a sample
func VarS(input []float64) float64 {
	return Var(input, 1)
}
