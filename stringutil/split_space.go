/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stringutil

import "unicode"

// SplitSpace splits a string on whitespace.
func SplitSpace(text string) (output []string) {
	if len(text) == 0 {
		return
	}

	var state int
	var word string
	for _, r := range text {
		switch state {
		case 0:	// word
			if unicode.IsSpace(r) {
				if len(word) > 0 {
					output = append(output, word)
					word = ""
				}
				state = 1
			} else {
				word = word + string(r)
			}
		case 1:
			if !unicode.IsSpace(r) {
				word = string(r)
				state = 0
			}
		}
	}

	if len(word) > 0 {
		output = append(output, word)
	}
	return
}
