package logger

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestLabels(t *testing.T) {
	assert := assert.New(t)

	labels := Labels{"foo": "bar", "fuzz": "buzz"}
	decomposed := labels.Decompose()
	assert.Len(decomposed, 2)
}
