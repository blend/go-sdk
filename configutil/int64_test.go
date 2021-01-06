package configutil

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestInt646464646464(t *testing.T) {
	assert := assert.New(t)

	intValue := Int64(0)
	ptr, err := intValue.Int64(context.TODO())
	assert.Nil(ptr)
	assert.Nil(err)

	intValue = Int64(1234)
	ptr, err = intValue.Int64(context.TODO())
	assert.Nil(err)
	assert.NotNil(ptr)
	assert.Equal(1234, *ptr)
}
