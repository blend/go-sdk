package validate

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestFirst(t *testing.T) {
	assert := assert.New(t)

	res := First(none, some(fmt.Errorf("one")), some(fmt.Errorf("two")), none)()
	assert.Equal(fmt.Errorf("one"), res)
}
