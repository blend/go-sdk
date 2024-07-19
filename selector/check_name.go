/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package selector

import (
	"unicode/utf8"
)

// CheckName checks the characters in a name but does not validate the length.
func CheckName(value string) (err error) {
	valueLen := len(value)
	var state int
	var ch rune
	var width int
	for pos := 0; pos < valueLen; pos += width {
		ch, width = utf8.DecodeRuneInString(value[pos:])
		switch state {
		case 0: //check prefix/suffix
			if !isAlpha(ch) {
				err = ErrLabelInvalidCharacter
				return
			}
			state = 1
			continue
		case 1:
			if !(isNameSymbol(ch) || isAlpha(ch)) {
				err = ErrLabelInvalidCharacter
				return
			}
			if pos == valueLen-2 {
				state = 0
			}
			continue
		}
	}
	return
}
