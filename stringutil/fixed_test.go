package stringutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestFixed(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("   abc", Fixed("abc", 6))
	assert.Equal("a", Fixed("abc", 1))
}

func TestFixedLeft(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("abc   ", FixedLeft("abc", 6))
	assert.Equal("a", FixedLeft("abc", 1))
}
