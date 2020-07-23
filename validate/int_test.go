package validate

import (
	"testing"

	"github.com/blend/go-sdk/assert"
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

	verr = Int(nil).Min(10)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrIntMin, ErrCause(verr))

	val = 1
	verr = Int(&val).Min(10)()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrIntMin, ErrCause(verr))

	val = 10
	verr = Int(&val).Min(10)()
	assert.Nil(verr)
}

func TestIntMax(t *testing.T) {
	assert := assert.New(t)

	var verr error
	val := 1
	verr = Int(&val).Max(10)()
	assert.Nil(verr)

	verr = Int(nil).Max(10)()
	assert.Nil(verr)

	val = 10
	verr = Int(&val).Max(10)()
	assert.Nil(verr)

	val = 11
	verr = Int(&val).Max(10)()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrIntMax, ErrCause(verr))
}

func TestIntBetween(t *testing.T) {
	assert := assert.New(t)

	var verr error
	val := 5
	verr = Int(&val).Between(1, 10)()
	assert.Nil(verr)

	verr = Int(nil).Between(5, 10)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrIntMin, ErrCause(verr))

	val = 1
	verr = Int(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(1, ErrValue(verr))
	assert.Equal(ErrIntMin, ErrCause(verr))

	val = 5
	verr = Int(&val).Between(5, 10)()
	assert.Nil(verr)

	val = 10
	verr = Int(&val).Between(5, 10)()
	assert.Nil(verr)

	val = 11
	verr = Int(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(11, ErrValue(verr))
	assert.Equal(ErrIntMax, ErrCause(verr))
}

func TestIntPositive(t *testing.T) {
	assert := assert.New(t)

	var verr error
	val := 5
	verr = Int(&val).Positive()()
	assert.Nil(verr)

	verr = Int(nil).Positive()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrIntPositive, ErrCause(verr))

	val = -5
	verr = Int(&val).Positive()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrIntPositive, ErrCause(verr))
}

func TestIntNegative(t *testing.T) {
	assert := assert.New(t)

	var verr error
	val := -5
	verr = Int(&val).Negative()()
	assert.Nil(verr)

	verr = Int(nil).Negative()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrIntNegative, ErrCause(verr))

	val = 5
	verr = Int(&val).Negative()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrIntNegative, ErrCause(verr))
}

func TestIntZero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	val := 0
	verr = Int(&val).Zero()()
	assert.Nil(verr)

	verr = Int(nil).Zero()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrIntZero, ErrCause(verr))

	val = 5
	verr = Int(&val).Zero()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrIntZero, ErrCause(verr))
}

func TestIntNotZero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	val := 5
	verr = Int(&val).NotZero()()
	assert.Nil(verr)

	verr = Int(nil).NotZero()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrIntNotZero, ErrCause(verr))

	val = 0
	verr = Int(&val).NotZero()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrIntNotZero, ErrCause(verr))
}
