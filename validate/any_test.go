package validate

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestAny(t *testing.T) {
	assert := assert.New(t)

	res := Any(nil, fmt.Errorf("one"), fmt.Errorf("two"), nil)
	assert.Equal(fmt.Errorf("one"), res)
}
