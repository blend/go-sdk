/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stringutil

const (
	// Empty is the empty string
	Empty	string	= ""

	// RuneSpace is a single rune representing a space.
	RuneSpace	rune	= ' '

	// RuneNewline is a single rune representing a newline.
	RuneNewline	rune	= '\n'

	// LowerA is the ascii int value for 'a'
	LowerA	uint	= uint('a')
	// LowerZ is the ascii int value for 'z'
	LowerZ	uint	= uint('z')
)

var (
	// LowerDiff is the difference between lower Z and lower A
	LowerDiff = (LowerZ - LowerA)
)
