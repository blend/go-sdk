/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package stringutil

// HasSuffixCaseless returns if a corpus has a suffix regardless of casing.
func HasSuffixCaseless(corpus, suffix string) bool {
	corpusLen := len(corpus)
	suffixLen := len(suffix)

	if corpusLen < suffixLen {
		return false
	}

	for x := 0; x < suffixLen; x++ {
		charCorpus := uint(corpus[corpusLen-(x+1)])
		charSuffix := uint(suffix[suffixLen-(x+1)])

		if charCorpus-LowerA <= LowerDiff {
			charCorpus = charCorpus - 0x20
		}

		if charSuffix-LowerA <= LowerDiff {
			charSuffix = charSuffix - 0x20
		}
		if charCorpus != charSuffix {
			return false
		}
	}
	return true
}
