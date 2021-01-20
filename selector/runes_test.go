package selector

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_isAlpha(t *testing.T) {
	its := assert.New(t)

	its.True(isAlpha('A'))
	its.True(isAlpha('a'))
	its.True(isAlpha('Z'))
	its.True(isAlpha('z'))
	its.True(isAlpha('0'))
	its.True(isAlpha('9'))
	its.True(isAlpha('함'))
	its.True(isAlpha('é'))
	its.False(isAlpha('-'))
	its.False(isAlpha('/'))
	its.False(isAlpha('~'))
}
