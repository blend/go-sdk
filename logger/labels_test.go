package logger

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestLabels(t *testing.T) {
	assert := assert.New(t)

	labels := make(Labels)
	labels.SetLabel("foo", "bar")
	labels.SetLabel("fuzz", "buzz")

	assert.Any(labels.GetLabelKeys(), func(v interface{}) bool {
		return v.(string) == "foo"
	})
	assert.Any(labels.GetLabelKeys(), func(v interface{}) bool {
		return v.(string) == "fuzz"
	})

	assert.Equal("bar", labels["foo"])
	assert.Equal("buzz", labels["fuzz"])

	value, ok := labels.GetLabel("foo")
	assert.True(ok)
	assert.Equal("bar", value)

	value, ok = labels.GetLabel("---")
	assert.False(ok)
	assert.Empty(value)

	value, ok = labels.GetLabel("fuzz")
	assert.True(ok)
	assert.Equal("buzz", value)

	decomposed := labels.Decompose()
	assert.Len(decomposed, 2)
}
