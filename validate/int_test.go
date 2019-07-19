package validate

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestIntMin(t *testing.T) {
	assert := assert.New(t)

	var verr error
	val := 10
	verr = Int(&val).Min(1)()
	assert.Nil(verr)

	val = 10
	verr = Int(&val).Min(10)()
	assert.Nil(verr)

	val = 1
	verr = Int(&val).Min(10)()
	assert.NotNil(verr)
	assert.Equal(ErrIntMin, ex.ErrInner(verr))
}

func TestIntMax(t *testing.T) {
	assert := assert.New(t)

	var verr error
	val := 1
	verr = Int(&val).Max(10)()
	assert.Nil(verr)

	val = 10
	verr = Int(&val).Max(10)()
	assert.Nil(verr)

	val = 10
	verr = Int(&val).Max(1)()
	assert.NotNil(verr)
	assert.Equal(ErrIntMax, ex.ErrInner(verr))
}

func TestIntBetween(t *testing.T) {
	assert := assert.New(t)

	var verr error
	val := 5
	verr = Int(&val).Between(1, 10)()
	assert.Nil(verr)

	val = 1
	verr = Int(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(ErrIntMin, ex.ErrInner(verr))

	val = 11
	verr = Int(&val).Between(1, 10)()
	assert.NotNil(verr)
	assert.Equal(ErrIntMax, ex.ErrInner(verr))
}

func TestIntPositive(t *testing.T) {
	assert := assert.New(t)

	var verr error
	val := 5
	verr = Int(&val).Positive()()
	assert.Nil(verr)

	val = -5
	verr = Int(&val).Positive()()
	assert.NotNil(verr)
	assert.Equal(ErrIntPositive, ex.ErrInner(verr))
}

func TestIntNegative(t *testing.T) {
	assert := assert.New(t)

	var verr error
	val := -5
	verr = Int(&val).Negative()()
	assert.Nil(verr)

	val = 5
	verr = Int(&val).Negative()()
	assert.NotNil(verr)
	assert.Equal(ErrIntNegative, ex.ErrInner(verr))
}
