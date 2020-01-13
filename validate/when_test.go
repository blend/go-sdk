package validate

import (
	"fmt"
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestWhen(t *testing.T) {
	assert := assert.New(t)

	var toggle bool
	when := When(func() bool { return toggle }, func() error { return fmt.Errorf("passes") })
	assert.Nil(when())
	toggle = true
	assert.Equal(fmt.Errorf("passes"), when())
}
