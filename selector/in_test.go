package selector

import (
	"testing"
	assert "github.com/blend/go-sdk/assert"
)

func TestIn(t *testing.T) {
	assert := assert.New(t)

	valid := Labels{
		"foo": "far",
		"moo": "lar",
	}
	valid2 := Labels{
		"foo": "bar",
		"moo": "lar",
	}
	missing := Labels{
		"loo": "mar",
		"moo": "lar",
	}
	invalid := Labels{
		"foo": "mar",
		"moo": "lar",
	}

	selector := In{Key: "foo", Values: []string{"bar", "far"}}
	assert.True(selector.Matches(valid))
	assert.True(selector.Matches(valid2))
	assert.True(selector.Matches(missing))
	assert.False(selector.Matches(invalid))

	assert.Equal("foo in (bar, far)", selector.String())
}
