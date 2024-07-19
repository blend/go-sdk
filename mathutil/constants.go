/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package mathutil

import "math"

const (
	_2pi = 2 * math.Pi
	_d2r = (math.Pi / 180.0)
	_r2d = (180.0 / math.Pi)

	// Epsilon represents the minimum amount of relevant delta we care about.
	Epsilon = 0.00000001
)
