package validate

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestIntMin(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Int.Min(1)(10)
	assert.Nil(verr)

	verr = Int.Min(10)(10)
	assert.Nil(verr)

	verr = Int.Min(10)(1)
	assert.NotNil(verr)
	assert.Equal(ErrIntMin, ex.ErrInner(verr))
}

func TestIntMax(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Int.Max(10)(1)
	assert.Nil(verr)

	verr = Int.Max(10)(10)
	assert.Nil(verr)

	verr = Int.Max(1)(10)
	assert.NotNil(verr)
	assert.Equal(ErrIntMax, ex.ErrInner(verr))
}

func TestIntBetween(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Int.Between(1, 10)(5)
	assert.Nil(verr)

	verr = Int.Between(5, 10)(1)
	assert.NotNil(verr)
	assert.Equal(ErrIntMin, ex.ErrInner(verr))

	verr = Int.Between(1, 10)(11)
	assert.NotNil(verr)
	assert.Equal(ErrIntMax, ex.ErrInner(verr))
}

func TestIntPositive(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Int.Positive(5)
	assert.Nil(verr)

	verr = Int.Positive(-5)
	assert.NotNil(verr)
	assert.Equal(ErrIntPositive, ex.ErrInner(verr))
}

func TestIntNegative(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Int.Negative(-5)
	assert.Nil(verr)

	verr = Int.Negative(5)
	assert.NotNil(verr)
	assert.Equal(ErrIntNegative, ex.ErrInner(verr))
}
