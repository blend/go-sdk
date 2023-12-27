/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestRewriteRuleApply(t *testing.T) {
	assert := assert.New(t)

	regex := `([0-9]+)\.([a-zA-Z]+)`
	expression := regexp.MustCompile(regex)

	rr := &RewriteRule{
		MatchExpression:	regex,
		expr:			expression,
		Action: func(path string, pieces ...string) string {
			assert.NotEmpty(path)
			assert.NotEmpty(pieces)
			assert.Len(pieces, 3, fmt.Sprintf("%#v", pieces))
			return path + "_ok!"
		},
	}

	matches, result := rr.Apply("1234.abcde")
	assert.True(matches)
	assert.Equal("1234.abcde_ok!", result)

	matches, result = rr.Apply("abcde.1234")
	assert.False(matches)
	assert.Equal("abcde.1234", result)
}
