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

func TestStringBetween(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String.Between(1, 3)("a")
	assert.Nil(verr)

	verr = String.Between(1, 3)("large")
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMax, ex.ErrInner(verr))

	verr = String.Between(1, 3)("")
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMin, ex.ErrInner(verr))
}

func TestStringMatches(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String.Matches("foo$")("a foo")
	assert.Nil(verr)

	verr = String.Matches("foo$")("foo not")
	assert.NotNil(verr)
	assert.Equal(ErrStringMatches, ex.ErrInner(verr))
}

func TestStringMatchesError(t *testing.T) {
	assert := assert.New(t)

	var err error
	err = String.Matches("((")("a foo") // this should be an invalid regex "(("
	assert.NotNil(err)
	assert.NotEqual(ErrValidation, ex.ErrClass(err))
}
