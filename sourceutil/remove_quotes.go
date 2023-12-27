/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package sourceutil

// RemoveQuotes removes quotes from a string
func RemoveQuotes(value string) string {
	var output []rune
	for _, r := range value {
		switch r {
		case '"', '\'':
			continue
		default:
			output = append(output, r)
		}
	}
	return string(output)
}
