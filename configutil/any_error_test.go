package configutil

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestAnyError(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(AnyError())
	assert.Nil(AnyError(nil))
	assert.Nil(AnyError(nil, nil))
	assert.Equal(fmt.Errorf("test"), AnyError(fmt.Errorf("test"), fmt.Errorf("test2")))
}
