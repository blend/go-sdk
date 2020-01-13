package stringutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestTrimSuffixCaseless(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("abc", TrimSuffixCaseless("abcdef", "def"))
	assert.Equal("ab2", TrimSuffixCaseless("ab2def", "DEF"))
	assert.Equal("ab3", TrimSuffixCaseless("ab3DEF", "def"))
	assert.Equal("abcdef", TrimSuffixCaseless("abcdef", "foo"))
	assert.Equal("abc", TrimSuffixCaseless("abc", "abcdef"))
}
