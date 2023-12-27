/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package selector

import "unicode"

func isWhitespace(ch rune) bool {
	return unicode.IsSpace(ch)
}

func isNameSymbol(ch rune) bool {
	switch ch {
	case Dot, Dash, Underscore:
		return true
	}
	return false
}

func isOperatorSymbol(ch rune) bool {
	return ch == Bang || ch == Equal
}

func isSymbol(ch rune) bool {
	return (int(ch) >= int(Bang) && int(ch) <= int(ForwardSlash)) ||
		(int(ch) >= int(Colon) && int(ch) <= int(At)) ||
		(int(ch) >= int(OpenBracket) && int(ch) <= int(BackTick)) ||
		(int(ch) >= int(OpenCurly) && int(ch) <= int(Tilde))
}

func isAlpha(ch rune) bool {
	return !isWhitespace(ch) && !unicode.IsControl(ch) && !isSymbol(ch)
}

func isDNSAlpha(ch rune) bool {
	if unicode.IsDigit(ch) {
		return true
	}
	if unicode.IsLetter(ch) {
		return unicode.IsLower(ch)
	}
	return false
}
