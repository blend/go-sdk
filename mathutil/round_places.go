/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package mathutil

import "math"

// RoundPlaces a float to a specific decimal place or precision
func RoundPlaces(input float64, places int) float64 {
	if math.IsNaN(input) {
		return 0.0
	}

	sign := 1.0
	if input < 0 {
		sign = -1
		input *= -1
	}

	rounded := float64(0)
	precision := math.Pow(10, float64(places))
	digit := input * precision
	_, decimal := math.Modf(digit)

	if decimal >= 0.5 {
		rounded = math.Ceil(digit)
	} else {
		rounded = math.Floor(digit)
	}

	return rounded / precision * sign
}
