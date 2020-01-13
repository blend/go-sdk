package logger

import (
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestCombineAnnotations(t *testing.T) {
	assert := assert.New(t)

	assert.Empty(CombineAnnotations(nil, nil, nil))
	combined := CombineAnnotations(Annotations{"foo": "bar"}, nil, Annotations{"moo": "loo"})
	assert.Equal("bar", combined["foo"])
	assert.Equal("loo", combined["moo"])
}
