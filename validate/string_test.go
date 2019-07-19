package validate

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestStringMin(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String("large").MinLen(3)()
	assert.Nil(verr)
	verr = String("a").MinLen(3)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMin, ex.ErrInner(verr))
}

func TestStringMaxlen(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String("a").MaxLen(3)()
	assert.Nil(verr)
	verr = String("large").MaxLen(3)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMax, ex.ErrInner(verr))
}

func TestStringBetweenLen(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String("a").BetweenLen(1, 3)()
	assert.Nil(verr)

	verr = String("large").BetweenLen(1, 3)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMax, ex.ErrInner(verr))

	verr = String("").BetweenLen(1, 3)()
	assert.NotNil(verr)
	assert.Equal(ErrValidation, ex.ErrClass(verr))
	assert.Equal(ErrStringLengthMin, ex.ErrInner(verr))
}

func TestStringMatches(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = String("a foo").Matches("foo$")()
	assert.Nil(verr)

	verr = String("foo not").Matches("foo$")()
	assert.NotNil(verr)
	assert.Equal(ErrStringMatches, ex.ErrInner(verr))
}

func TestStringMatchesError(t *testing.T) {
	assert := assert.New(t)

	var err error
	err = String("a foo").Matches("((")() // this should be an invalid regex "(("
	assert.NotNil(err)
	assert.NotEqual(ErrValidation, ex.ErrClass(err))
}
