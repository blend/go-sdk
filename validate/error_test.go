package validate

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestError(t *testing.T) {
	assert := assert.New(t)

	verr := Error(fmt.Errorf("this is a test"), nil)
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Empty(ex.ErrMessage(verr))
	assert.Equal(fmt.Errorf("this is a test"), Cause(verr))

	verr = Error(fmt.Errorf("this is a test"), nil, "foo", "bar")
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal("foobar", Message(verr))
	assert.Equal(fmt.Errorf("this is a test"), Cause(verr))
}

func TestErrorf(t *testing.T) {
	assert := assert.New(t)

	verr := Errorf(fmt.Errorf("this is a test"), "foo", "minimum: %d", 30)
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal("minimum: 30", Message(verr))
	assert.Equal(fmt.Errorf("this is a test"), Cause(verr))
}

func TestCause(t *testing.T) {
	assert := assert.New(t)

	err := ex.New(ErrNonLengthType)
	assert.Equal(ErrNonLengthType, ex.ErrClass(err))
	assert.Equal(ErrNonLengthType, Cause(err))

	verr := Error(ErrEmpty, "foo")
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrEmpty, Cause(verr))
}
