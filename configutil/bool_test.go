package configutil

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestBool(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(Bool(nil))
	ret, err := (*BoolValue)(nil).Bool(context.TODO())
	assert.Nil(ret)
	assert.Nil(err)

	value := true
	bv := Bool(&value)
	assert.NotNil(bv)

	ret, err = bv.Bool(context.TODO())
	assert.Nil(err)
	assert.NotNil(ret)
	assert.True(*ret)
}
