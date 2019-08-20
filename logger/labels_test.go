package logger

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCombineLabels(t *testing.T) {
	assert := assert.New(t)

	combined := CombineLabels(Labels{"foo": "bar"}, Labels{"moo": "loo"})
	assert.Equal("bar", combined["foo"])
	assert.Equal("loo", combined["moo"])
}
