package stringutil

import (
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestHasSuffixCaseless(t *testing.T) {
	assert := assert.New(t)

	assert.True(HasSuffixCaseless("hello world!", "world!"))
	assert.True(HasSuffixCaseless("hello world", "d"))
	assert.True(HasSuffixCaseless("hello world", "hello world"))

	assert.True(HasSuffixCaseless("hello WORLD", "world"))
	assert.True(HasSuffixCaseless("hello world", "WORLD"))

	assert.False(HasSuffixCaseless("hello world", "hello hello world"))
	assert.False(HasSuffixCaseless("hello world", "foobar"))
	assert.False(HasSuffixCaseless("hello world", "hello")) //this would pass prefix
}
