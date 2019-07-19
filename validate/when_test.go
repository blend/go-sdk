package validate

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestWhen(t *testing.T) {
	assert := assert.New(t)

	var toggle bool
	when := When(func() bool { return toggle }, func() error { return fmt.Errorf("passes") }, func() error { return fmt.Errorf("fails") })

	err := when()
	assert.Equal(fmt.Errorf("fails"), err)

	toggle = true

	err = when()
	assert.Equal(fmt.Errorf("passes"), err)
}
