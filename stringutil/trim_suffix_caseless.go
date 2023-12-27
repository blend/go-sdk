/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stringutil

// TrimSuffixCaseless trims a case insensitive suffix from a corpus.
func TrimSuffixCaseless(corpus, suffix string) string {
	corpusLen := len(corpus)
	suffixLen := len(suffix)

	if corpusLen < suffixLen {
		return corpus
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
			return corpus
		}
	}
	return corpus[:corpusLen-suffixLen]
}
