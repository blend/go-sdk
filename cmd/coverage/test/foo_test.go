package foo

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestBar(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(1, bar())
}
