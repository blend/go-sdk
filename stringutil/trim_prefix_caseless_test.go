package stringutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestTrimPrefixCaseless(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("def", TrimPrefixCaseless("abcdef", "abc"))
	assert.Equal("def", TrimPrefixCaseless("abcdef", "ABC"))
	assert.Equal("DEF", TrimPrefixCaseless("abcDEF", "abc"))
	assert.Equal("abcdef", TrimPrefixCaseless("abcdef", "foo"))
	assert.Equal("abc", TrimPrefixCaseless("abc", "abcdef"))
}
