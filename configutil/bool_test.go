package configutil

import (
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestBool(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(Bool(nil))
	ret, err := (*BoolValue)(nil).Bool()
	assert.Nil(ret)
	assert.Nil(err)

	value := true
	bv := Bool(&value)
	assert.NotNil(bv)

	ret, err = bv.Bool()
	assert.Nil(err)
	assert.NotNil(ret)
	assert.True(*ret)
}
