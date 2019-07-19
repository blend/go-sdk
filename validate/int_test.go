package validate

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestIntMin(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Int(10).Min(1)()
	assert.Nil(verr)

	verr = Int(10).Min(10)()
	assert.Nil(verr)

	verr = Int(1).Min(10)()
	assert.NotNil(verr)
	assert.Equal(ErrIntMin, ex.ErrInner(verr))
}

func TestIntMax(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Int(1).Max(10)()
	assert.Nil(verr)

	verr = Int(10).Max(10)()
	assert.Nil(verr)

	verr = Int(10).Max(1)()
	assert.NotNil(verr)
	assert.Equal(ErrIntMax, ex.ErrInner(verr))
}

func TestIntBetween(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Int(5).Between(1, 10)()
	assert.Nil(verr)

	verr = Int(1).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(ErrIntMin, ex.ErrInner(verr))

	verr = Int(11).Between(1, 10)()
	assert.NotNil(verr)
	assert.Equal(ErrIntMax, ex.ErrInner(verr))
}

func TestIntPositive(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Int(5).Positive()
	assert.Nil(verr)

	verr = Int(-5).Positive()
	assert.NotNil(verr)
	assert.Equal(ErrIntPositive, ex.ErrInner(verr))
}

func TestIntNegative(t *testing.T) {
	assert := assert.New(t)

	var verr error
	verr = Int(-5).Negative()
	assert.Nil(verr)

	verr = Int(5).Negative()
	assert.NotNil(verr)
	assert.Equal(ErrIntNegative, ex.ErrInner(verr))
}
