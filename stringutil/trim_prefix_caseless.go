/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package stringutil

// TrimPrefixCaseless trims a prefix from a corpus ignoring case.
func TrimPrefixCaseless(corpus, prefix string) string {
	corpusLen := len(corpus)
	prefixLen := len(prefix)

	if corpusLen < prefixLen {
		return corpus
	}

	for x := 0; x < prefixLen; x++ {
		charCorpus := uint(corpus[x])
		charPrefix := uint(prefix[x])

		if charCorpus-LowerA <= LowerDiff {
			charCorpus = charCorpus - 0x20
		}

		if charPrefix-LowerA <= LowerDiff {
			charPrefix = charPrefix - 0x20
		}
		if charCorpus != charPrefix {
			return corpus
		}
	}

	return corpus[prefixLen:]
}
