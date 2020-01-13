package configutil

import (
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestInt(t *testing.T) {
	assert := assert.New(t)

	intValue := Int(1234)
	ptr, err := intValue.Int()
	assert.Nil(err)
	assert.NotNil(ptr)
	assert.Equal(1234, *ptr)
}
