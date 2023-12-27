/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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
		{corpus: "foo/${bar}/baz", expected: "foo/example-string/baz", tokens: map[string]string{"bar": "example-string"}, message: "should handle escaping a single variable"},
		{corpus: "foo/${what}/baz", expected: "foo/${what}/baz", tokens: map[string]string{"bar": "example-string"}, message: "should handle unknown variables"},
		{corpus: "foo/${bar}/baz/${buzz}", expected: "foo/example-string/baz/dog", tokens: map[string]string{"bar": "example-string", "buzz": "dog"}, message: "should handle escaping multiple variables"},
		{corpus: "foo/${bar${buzz}foo}/bar", expected: "foo/${bar${buzz}foo}/bar", tokens: map[string]string{"bar": "example-string", "buzz": "dog"}, message: "nesting variables should produce a weird key"},
	}

	for _, testCase := range testCases {
		assert.Equal(testCase.expected, Tokenize(testCase.corpus, testCase.tokens), testCase.message)
	}
}
