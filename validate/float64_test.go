/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package validate

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestFloat64Min(t *testing.T) {
	assert := assert.New(t)

	var verr error
	val := 10.0
	verr = Float64(&val).Min(1)()
	assert.Nil(verr)

	val = 10.0
	verr = Float64(&val).Min(10)()
	assert.Nil(verr)

	verr = Float64(nil).Min(10)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrFloat64Min, ErrCause(verr))

	val = 1.0
	verr = Float64(&val).Min(10)()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrFloat64Min, ErrCause(verr))

	val = 10.0
	verr = Float64(&val).Min(10)()
	assert.Nil(verr)
}

func TestFloat64Max(t *testing.T) {
	assert := assert.New(t)

	var verr error
	val := 1.0
	verr = Float64(&val).Max(10)()
	assert.Nil(verr)

	verr = Float64(nil).Max(10)()
	assert.Nil(verr)

	val = 10.0
	verr = Float64(&val).Max(10)()
	assert.Nil(verr)

	val = 11.0
	verr = Float64(&val).Max(10)()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrFloat64Max, ErrCause(verr))
}

func TestFloat64Between(t *testing.T) {
	assert := assert.New(t)

	var verr error
	val := 5.0
	verr = Float64(&val).Between(1, 10)()
	assert.Nil(verr)

	verr = Float64(nil).Between(5, 10)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrFloat64Min, ErrCause(verr))

	val = 1.0
	verr = Float64(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(1, ErrValue(verr))
	assert.Equal(ErrFloat64Min, ErrCause(verr))

	val = 5.0
	verr = Float64(&val).Between(5, 10)()
	assert.Nil(verr)

	val = 10.0
	verr = Float64(&val).Between(5, 10)()
	assert.Nil(verr)

	val = 11.0
	verr = Float64(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(11, ErrValue(verr))
	assert.Equal(ErrFloat64Max, ErrCause(verr))
}

func TestFloat64Positive(t *testing.T) {
	assert := assert.New(t)

	var verr error
	val := 5.0
	verr = Float64(&val).Positive()()
	assert.Nil(verr)

	verr = Float64(nil).Positive()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrFloat64Positive, ErrCause(verr))

	val = -5.0
	verr = Float64(&val).Positive()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrFloat64Positive, ErrCause(verr))
}

func TestFloat64Negative(t *testing.T) {
	assert := assert.New(t)

	var verr error
	val := -5.0
	verr = Float64(&val).Negative()()
	assert.Nil(verr)

	verr = Float64(nil).Negative()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrFloat64Negative, ErrCause(verr))

	val = 5.0
	verr = Float64(&val).Negative()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrFloat64Negative, ErrCause(verr))
}

func TestFloat64Epsilon(t *testing.T) {
	assert := assert.New(t)

	var verr error
	val := 5.0
	verr = Float64(&val).Epsilon(4.999999, DefaultEpsilon)()
	assert.Nil(verr)

	verr = Float64(nil).Epsilon(4.999999, DefaultEpsilon)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrFloat64Epsilon, ErrCause(verr))

	verr = Float64(&val).Epsilon(4.99, DefaultEpsilon)()
	assert.NotNil(verr)
	assert.Equal(5.0, ErrValue(verr))
	assert.Equal(ErrFloat64Epsilon, ErrCause(verr))
}

func TestFloat64Zero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	val := 0.0
	verr = Float64(&val).Zero()()
	assert.Nil(verr)

	verr = Float64(nil).Zero()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrFloat64Zero, ErrCause(verr))

	val = 5.0
	verr = Float64(&val).Zero()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrFloat64Zero, ErrCause(verr))
}

func TestFloat64NotZero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	val := 5.0
	verr = Float64(&val).NotZero()()
	assert.Nil(verr)

	verr = Float64(nil).NotZero()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrFloat64NotZero, ErrCause(verr))

	val = 0.0
	verr = Float64(&val).NotZero()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrFloat64NotZero, ErrCause(verr))
}
