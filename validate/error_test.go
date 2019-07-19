package validate

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestError(t *testing.T) {
	assert := assert.New(t)

	verr := Error(fmt.Errorf("this is a test"))
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Empty(ex.ErrMessage(verr))
	assert.Equal(fmt.Errorf("this is a test"), ex.ErrInner(verr))

	verr = Error(fmt.Errorf("this is a test"), "foo", "bar")
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal("foobar", ex.ErrMessage(verr))
	assert.Equal(fmt.Errorf("this is a test"), ex.ErrInner(verr))
}

func TestErrorf(t *testing.T) {
	assert := assert.New(t)

	verr := Errorf(fmt.Errorf("this is a test"), "minimum: %d", 30)
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal("minimum: 30", ex.ErrMessage(verr))
	assert.Equal(fmt.Errorf("this is a test"), ex.ErrInner(verr))
}
