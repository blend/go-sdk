/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package mathutil

import "math"

// RoundUp rounds up to a given roundTo value.
func RoundUp(value, roundTo float64) float64 {
	d1 := math.Ceil(value / roundTo)
	return d1 * roundTo
}

// RoundDown rounds down to a given roundTo value.
func RoundDown(value, roundTo float64) float64 {
	d1 := math.Floor(value / roundTo)
	return d1 * roundTo
}
