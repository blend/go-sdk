/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package slant

// Font represents a single font.
type Font struct {
	// Height of one char
	Height int
	// Baseline is the height of letters not including descenders.
	Baseline int
	// Width of the widest char
	Width int
	// Hardblank symbol is the non-smushable space character.
	Hardblank rune
	// A string for each line of the char
	Letters [][]string
}
