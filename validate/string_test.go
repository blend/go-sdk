package validate

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestStringMin(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String.Min(3)("large")
	assert.Nil(verr)
	verr = String.Min(3)("a")
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMin, ex.ErrInner(verr))
}

func TestStringMax(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String.Max(3)("a")
	assert.Nil(verr)
	verr = String.Max(3)("large")
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMax, ex.ErrInner(verr))
}
