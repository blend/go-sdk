package stringutil

import (
	"unicode"
)

// Slugify replaces non-letter or digit runes with '-'.
func Slugify(v string) string {
	runes := []rune(v)
	var output []rune
	var c rune
	var previousWasReplaced bool
	for index := range runes {
		c = runes[index]
		if unicode.IsLetter(c) {
			output = append(output, unicode.ToLower(c))
			previousWasReplaced = false
			continue
		}
		if unicode.IsDigit(c) {
			output = append(output, c)
			previousWasReplaced = false
			continue
		}
		if unicode.IsSpace(c) {
			if !previousWasReplaced {
				output = append(output, '-')
				previousWasReplaced = true
			}
			continue
		}
		previousWasReplaced = false
	}
	return string(output)
}
