package logger

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCombineAnnotations(t *testing.T) {
	assert := assert.New(t)

	combined := CombineAnnotations(Annotations{"foo": "bar"}, Annotations{"moo": "loo"})
	assert.Equal("bar", combined["foo"])
	assert.Equal("loo", combined["moo"])
}
