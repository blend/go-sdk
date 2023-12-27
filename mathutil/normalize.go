/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package mathutil

// Normalize returns a set of numbers on the interval [0,1] for a given set of inputs.
// An example: 4,3,2,1 => 0.4, 0.3, 0.2, 0.1
// Caveat; the total may be < 1.0; there are going to be issues with irrational numbers etc.
func Normalize(values ...float64) []float64 {
	var total float64
	for _, v := range values {
		total += v
	}
	output := make([]float64, len(values))
	for x, v := range values {
		output[x] = RoundDown(v/total, 0.0001)
	}
	return output
}
