/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stringutil

import "sort"

var (
	// LowerLetters is a runset of lowercase letters.
	LowerLetters Runeset = []rune("abcdefghijklmnopqrstuvwxyz")

	// UpperLetters is a runset of uppercase letters.
	UpperLetters Runeset = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	// Letters is a runset of both lower and uppercase letters.
	Letters = append(LowerLetters, UpperLetters...)

	// Numbers is a runset of numeric characters.
	Numbers Runeset = []rune("0123456789")

	// LettersAndNumbers is a runset of letters and numeric characters.
	LettersAndNumbers = append(Letters, Numbers...)

	// Symbols is a runset of symbol characters.
	Symbols Runeset = []rune(`!@#$%^&*()_+-=[]{}\|:;`)

	// LettersNumbersAndSymbols is a runset of letters, numbers and symbols.
	LettersNumbersAndSymbols = append(LettersAndNumbers, Symbols...)
)

// CombineRunsets combines given runsets into a single runset.
func CombineRunsets(runesets ...[]rune) []rune {
	output := []rune{}
	for _, set := range runesets {
		output = append(output, set...)
	}
	return output
}

// Runeset is a set of runes
type Runeset []rune

// Len implements part of sorter.
func (rs Runeset) Len() int {
	return len(rs)
}

// Swap implements part of sorter.
func (rs Runeset) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

// Less implements part of sorter.
func (rs Runeset) Less(i, j int) bool {
	return uint16(rs[i]) < uint16(rs[j])
}

// Set returns a map of the runes in the set.
func (rs Runeset) Set() map[rune]bool {
	seen := make(map[rune]bool)
	for _, r := range rs {
		seen[r] = true
	}
	return seen
}

// Combine merges runesets.
func (rs Runeset) Combine(other ...Runeset) Runeset {
	seen := rs.Set()
	for _, set := range other {
		for r := range set.Set() {
			seen[r] = true
		}
	}

	var output []rune
	for r := range seen {
		output = append(output, r)
	}

	sort.Sort(Runeset(output))
	return Runeset(output)
}

// Random returns a random selection of runes from the set.
func (rs Runeset) Random(length int) string {
	runes := make([]rune, length)
	for index := range runes {
		runes[index] = rs[provider.Intn(len(rs)-1)]
	}
	return string(runes)
}
