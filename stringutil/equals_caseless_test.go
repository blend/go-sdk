package stringutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestEqualsCaseless(t *testing.T) {
	assert := assert.New(t)
	assert.True(EqualsCaseless("foo", "FOO"))
	assert.True(EqualsCaseless("foo123", "FOO123"))
	assert.True(EqualsCaseless("!foo123", "!foo123"))
	assert.False(EqualsCaseless("foo", "bar"))
}
