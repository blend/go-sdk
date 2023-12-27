/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stringutil

// ReplaceAny replaces any runes in the 'replaced' list with a given replacement.
// Example:
//
//	output := ReplaceAny("foo bar_baz", '-', []rune(` _`)...)
func ReplaceAny(corpus string, replacement rune, replaced ...rune) string {
	characters := []rune(corpus)
	var c rune
	for x := 0; x < len(characters); x++ {
		c = characters[x]
		for _, r := range replaced {
			if c == r {
				characters[x] = replacement
				break
			}
		}
	}

	return string(characters)
}
