package selector

import (
	"testing"

	assert "github.com/blend/go-sdk/assert"
)

func TestNotHasKey(t *testing.T) {
	assert := assert.New(t)

	valid := Labels{
		"foo": "far",
	}
	assert.False(NotHasKey("foo").Matches(valid))
	assert.True(NotHasKey("zoo").Matches(valid))
	assert.Equal("!foo", NotHasKey("foo").String())
}
