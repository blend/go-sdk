package stringutil

import (
	"unicode"
)

// Slugify replaces non-letter or digit runes with '-'.
// It will not add repeated '-'.
func Slugify(v string) string {
	runes := []rune(v)

	var output []rune
	var c rune
	var lastWasDash bool

	for index := range runes {
		c = runes[index]
		// add letters and make them lowercase
		if unicode.IsLetter(c) {
			output = append(output, unicode.ToLower(c))
			lastWasDash = false
			continue
		}
		// add digits unchanged
		if unicode.IsDigit(c) {
			output = append(output, c)
			lastWasDash = false
			continue
		}
		// if we hit a dash, only add it if
		// the last character wasnt a dash
		if c == '-' {
			if !lastWasDash {
				output = append(output, c)
				lastWasDash = true
			}
			continue

		}
		if unicode.IsSpace(c) {
			if !lastWasDash {
				output = append(output, '-')
				lastWasDash = true
			}
			continue
		}
	}
	return string(output)
}
