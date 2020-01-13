package stringutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

type replaceAnyTestCase struct {
	expected  string
	corpus    string
	with      rune
	toReplace []rune
}

func TestReplaceAny(t *testing.T) {
	assert := assert.New(t)

	testCases := []replaceAnyTestCase{
		{expected: "", corpus: "", with: '_', toReplace: Symbols},
		{expected: "foo", corpus: "foo", with: '_', toReplace: Symbols},
		{expected: "foo_", corpus: "foo$", with: '_', toReplace: Symbols},
		{expected: "_foo_", corpus: "&foo$", with: '_', toReplace: Symbols},
		{expected: "_fo o_", corpus: "&fo o$", with: '_', toReplace: Symbols},
	}

	for _, testCase := range testCases {
		assert.Equal(testCase.expected, ReplaceAny(testCase.corpus, testCase.with, testCase.toReplace...))
	}
}
