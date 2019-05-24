package logger

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestLabels(t *testing.T) {
	assert := assert.New(t)

	labels := make(Labels)
	labels.AddLabelValue("foo", "bar")
	labels.AddLabelValue("fuzz", "buzz")

	assert.Equal("bar", labels["foo"])
	assert.Equal("buzz", labels["fuzz"])

	value, ok := labels.GetLabelValue("foo")
	assert.True(ok)
	assert.Equal("bar", value)

	value, ok = labels.GetLabelValue("---")
	assert.False(ok)
	assert.Empty(value)

	value, ok = labels.GetLabelValue("fuzz")
	assert.True(ok)
	assert.Equal("buzz", value)

	decomposed := labels.Decompose()
	assert.Len(decomposed, 2)
}
