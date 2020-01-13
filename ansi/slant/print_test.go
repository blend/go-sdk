package slant

import (
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestPrint(t *testing.T) {
	assert := assert.New(t)

	output, err := PrintString("WARDEN")
	assert.Nil(err)
	assert.NotEmpty(output)
}
