/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package mathutil

import "math"

// StdDevP finds the amount of variation from the population
func StdDevP(input []float64) float64 {
	if len(input) == 0 {
		return 0
	}

	// stdev is generally the square root of the variance
	return math.Pow(VarP(input), 0.5)
}

// StdDevS finds the amount of variation from a sample
func StdDevS(input []float64) float64 {
	if len(input) == 0 {
		return 0
	}

	// stdev is generally the square root of the variance
	return math.Pow(VarS(input), 0.5)
}
