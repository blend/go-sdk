package validate

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestFloat32Min(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val float32 = 10.0
	verr = Float32(&val).Min(1)()
	assert.Nil(verr)

	val = 10.0
	verr = Float32(&val).Min(10)()
	assert.Nil(verr)

	verr = Float32(nil).Min(10)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrFloat32Min, ErrCause(verr))

	val = 1.0
	verr = Float32(&val).Min(10)()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrFloat32Min, ErrCause(verr))

	val = 10.0
	verr = Float32(&val).Min(10)()
	assert.Nil(verr)
}

func TestFloat32Max(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val float32 = 1.0
	verr = Float32(&val).Max(10)()
	assert.Nil(verr)

	verr = Float32(nil).Max(10)()
	assert.Nil(verr)

	val = 10.0
	verr = Float32(&val).Max(10)()
	assert.Nil(verr)

	val = 11.0
	verr = Float32(&val).Max(10)()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrFloat32Max, ErrCause(verr))
}

func TestFloat32Between(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val float32 = 5.0
	verr = Float32(&val).Between(1, 10)()
	assert.Nil(verr)

	verr = Float32(nil).Between(5, 10)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrFloat32Min, ErrCause(verr))

	val = 1.0
	verr = Float32(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(1, ErrValue(verr))
	assert.Equal(ErrFloat32Min, ErrCause(verr))

	val = 5.0
	verr = Float32(&val).Between(5, 10)()
	assert.Nil(verr)

	val = 10.0
	verr = Float32(&val).Between(5, 10)()
	assert.Nil(verr)

	val = 11.0
	verr = Float32(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(11, ErrValue(verr))
	assert.Equal(ErrFloat32Max, ErrCause(verr))
}

func TestFloat32Positive(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val float32 = 5.0
	verr = Float32(&val).Positive()()
	assert.Nil(verr)

	verr = Float32(nil).Positive()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrFloat32Positive, ErrCause(verr))

	val = -5.0
	verr = Float32(&val).Positive()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrFloat32Positive, ErrCause(verr))
}

func TestFloat32Negative(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val float32 = -5.0
	verr = Float32(&val).Negative()()
	assert.Nil(verr)

	verr = Float32(nil).Negative()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrFloat32Negative, ErrCause(verr))

	val = 5.0
	verr = Float32(&val).Negative()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrFloat32Negative, ErrCause(verr))
}

func TestFloat32Epsilon(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val float32 = 5.0
	verr = Float32(&val).Epsilon(4.999999, DefaultEpsilon)()
	assert.Nil(verr)

	verr = Float32(nil).Epsilon(4.999999, DefaultEpsilon)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrFloat32Epsilon, ErrCause(verr))

	verr = Float32(&val).Epsilon(4.99, DefaultEpsilon)()
	assert.NotNil(verr)
	assert.Equal(5.0, ErrValue(verr))
	assert.Equal(ErrFloat32Epsilon, ErrCause(verr))
}

func TestFloat32Zero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val float32 = 0.0
	verr = Float32(&val).Zero()()
	assert.Nil(verr)

	verr = Float32(nil).Zero()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrFloat32Zero, ErrCause(verr))

	val = 5.0
	verr = Float32(&val).Zero()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrFloat32Zero, ErrCause(verr))
}

func TestFloat32NotZero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val float32 = 5.0
	verr = Float32(&val).NotZero()()
	assert.Nil(verr)

	verr = Float32(nil).NotZero()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrFloat32NotZero, ErrCause(verr))

	val = 0.0
	verr = Float32(&val).NotZero()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrFloat32NotZero, ErrCause(verr))
}
