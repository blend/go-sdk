/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package mathutil

import "math"

// PowInt returns the base to the power.
func PowInt(base int, power uint) int {
	if base == 2 {
		return 1 << power
	}
	return int(math.RoundToEven((math.Pow(float64(base), float64(power)))))
}
