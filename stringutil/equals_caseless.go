/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stringutil

// EqualsCaseless compares two strings regardless of case.
func EqualsCaseless(a, b string) bool {
	aLen := len(a)
	bLen := len(b)
	if aLen != bLen {
		return false
	}

	for x := 0; x < aLen; x++ {
		charA := uint(a[x])
		charB := uint(b[x])

		if charA-LowerA <= LowerDiff {
			charA = charA - 0x20
		}
		if charB-LowerA <= LowerDiff {
			charB = charB - 0x20
		}
		if charA != charB {
			return false
		}
	}

	return true
}
