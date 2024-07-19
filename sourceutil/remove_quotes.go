/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
