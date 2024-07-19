/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stats

import (
	"strings"
	"unicode"
)

// Tag formats a tag with a given key and value.
// For tags in the form `key` use an empty string for the value.
func Tag(key, value string) string {
	key = cleanTagElement(key)
	value = cleanTagElement(value)
	return key + ":" + value
}

// SplitTag splits a given tag in a key and a value
func SplitTag(tag string) (key, value string) {
	parts := strings.SplitN(tag, ":", 2)
	if len(parts) > 0 {
		key = parts[0]
	}
	if len(parts) > 1 {
		value = parts[1]
	}
	return
}

// cleansTagElement cleans up tag elements as best as it can
// per the spec at https://docs.datadoghq.com/tagging/
func cleanTagElement(value string) string {
	valueRunes := []rune(value)
	var r rune
	for x := 0; x < len(valueRunes); x++ {
		r = valueRunes[x]
		// letters
		if unicode.IsLetter(r) {
			continue
		}
		// digits
		if unicode.IsDigit(r) {
			continue
		}
		// allowed symbols
		switch r {
		case '-', ':', '_', '.', '/', '\\':
			continue
		default:
		}
		// everything else
		valueRunes[x] = '_'
		continue
	}
	return string(valueRunes)
}
