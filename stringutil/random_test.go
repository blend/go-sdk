package stringutil

import (
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestRandom(t *testing.T) {
	assert := assert.New(t)

	output := Random(Letters, 32)
	assert.Len(output, 32)

	output2 := Random(Letters, 32)
	assert.Len(output2, 32)

	assert.NotEqual(output, output2)
}
