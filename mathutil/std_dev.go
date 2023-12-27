/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
