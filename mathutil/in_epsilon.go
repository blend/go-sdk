/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package mathutil

// InEpsilon returns if two values are within the Epsilon of each other absolutely.
func InEpsilon(a, b float64) bool {
	return (a-b) < Epsilon && (b-a) < Epsilon
}
