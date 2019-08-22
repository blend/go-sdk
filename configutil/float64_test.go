package configutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestFloat64(t *testing.T) {
	assert := assert.New(t)

	floatValue := Float64(3.14)
	ptr, err := floatValue.Float64()
	assert.Nil(err)
	assert.NotNil(ptr)
	assert.Equal(3.14, *ptr)
}
