package stringutil

import (
	"bytes"
	"unicode"
)

// Title returns a string in title case.
func Title(corpus string) string {
	output := bytes.NewBuffer(nil)
	runes := []rune(corpus)

	haveSeenLetter := false
	var r rune
	for x := 0; x < len(runes); x++ {
		r = runes[x]

		if unicode.IsLetter(r) {
			if !haveSeenLetter {
				output.WriteRune(unicode.ToUpper(r))
				haveSeenLetter = true
			} else {
				output.WriteRune(unicode.ToLower(r))
			}
		} else {
			output.WriteRune(r)
			haveSeenLetter = false
		}
	}
	return output.String()
}
