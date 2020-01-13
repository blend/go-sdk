package stringutil

import (
	"testing"
	"github.com/blend/go-sdk/assert"
)

type tokenizeTestCase struct {
	corpus   string
	tokens   map[string]string
	expected string
	message  string
}

func TestStringTokenize(t *testing.T) {
	assert := assert.New(t)

	testCases := []tokenizeTestCase{
		{corpus: "", expected: "", message: "should handle the empty input case"},
		{corpus: "ff", expected: "ff", message: "should handle the (nearly) empty input case"},
		{corpus: "foo/${bar}/baz", expected: "foo/bailey/baz", tokens: map[string]string{"bar": "bailey"}, message: "should handle escaping a single variable"},
		{corpus: "foo/${what}/baz", expected: "foo/${what}/baz", tokens: map[string]string{"bar": "bailey"}, message: "should handle unknown variables"},
		{corpus: "foo/${bar}/baz/${buzz}", expected: "foo/bailey/baz/dog", tokens: map[string]string{"bar": "bailey", "buzz": "dog"}, message: "should handle escaping multiple variables"},
		{corpus: "foo/${bar${buzz}foo}/bar", expected: "foo/${bar${buzz}foo}/bar", tokens: map[string]string{"bar": "bailey", "buzz": "dog"}, message: "nesting variables should produce a weird key"},
	}

	for _, testCase := range testCases {
		assert.Equal(testCase.expected, Tokenize(testCase.corpus, testCase.tokens), testCase.message)
	}
}
