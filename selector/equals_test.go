package selector

import (
	"testing"

	assert "github.com/blend/go-sdk/assert"
)

func TestEquals(t *testing.T) {
	assert := assert.New(t)

	valid := Labels{
		"foo": "far",
		"moo": "bar",
	}
	assert.True(Equals{Key: "foo", Value: "far"}.Matches(valid))
	assert.False(Equals{Key: "zoo", Value: "buzz"}.Matches(valid))
	assert.False(Equals{Key: "foo", Value: "bar"}.Matches(valid))

	assert.Equal("foo == bar", Equals{Key: "foo", Value: "bar"}.String())

	// Test: selector option does not mutate the operator
	eq := Equals{Key: "foo", Value: "*far"}
	assert.Nil(eq.Validate(SelectorOptPermittedValues('*')))
	assert.Empty(eq.PermittedValues)
}
