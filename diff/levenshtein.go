/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package diff

import "unicode/utf8"

// Levenshtein computes the Levenshtein distance that is the number of inserted, deleted or substituted characters.
func Levenshtein(diffs []Diff) int {
	levenshtein := 0
	insertions := 0
	deletions := 0

	for _, aDiff := range diffs {
		switch aDiff.Type {
		case DiffInsert:
			insertions += utf8.RuneCountInString(aDiff.Text)
		case DiffDelete:
			deletions += utf8.RuneCountInString(aDiff.Text)
		case DiffEqual:
			// A deletion and an insertion is one substitution.
			levenshtein += max(insertions, deletions)
			insertions = 0
			deletions = 0
		}
	}

	levenshtein += max(insertions, deletions)
	return levenshtein
}
