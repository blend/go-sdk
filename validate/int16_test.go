/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package validate

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestInt16Min(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int16 = 10
	verr = Int16(&val).Min(1)()
	assert.Nil(verr)

	val = 10
	verr = Int16(&val).Min(10)()
	assert.Nil(verr)

	verr = Int16(nil).Min(10)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt16Min, ErrCause(verr))

	val = 1
	verr = Int16(&val).Min(10)()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt16Min, ErrCause(verr))

	val = 10
	verr = Int16(&val).Min(10)()
	assert.Nil(verr)
}

func TestInt16Max(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int16 = 1
	verr = Int16(&val).Max(10)()
	assert.Nil(verr)

	verr = Int16(nil).Max(10)()
	assert.Nil(verr)

	val = 10
	verr = Int16(&val).Max(10)()
	assert.Nil(verr)

	val = 11
	verr = Int16(&val).Max(10)()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt16Max, ErrCause(verr))
}

func TestInt16Between(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int16 = 5
	verr = Int16(&val).Between(1, 10)()
	assert.Nil(verr)

	verr = Int16(nil).Between(5, 10)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt16Min, ErrCause(verr))

	val = 1
	verr = Int16(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(1, ErrValue(verr))
	assert.Equal(ErrInt16Min, ErrCause(verr))

	val = 5
	verr = Int16(&val).Between(5, 10)()
	assert.Nil(verr)

	val = 10
	verr = Int16(&val).Between(5, 10)()
	assert.Nil(verr)

	val = 11
	verr = Int16(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(11, ErrValue(verr))
	assert.Equal(ErrInt16Max, ErrCause(verr))
}

func TestInt16Positive(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int16 = 5
	verr = Int16(&val).Positive()()
	assert.Nil(verr)

	verr = Int16(nil).Positive()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt16Positive, ErrCause(verr))

	val = -5
	verr = Int16(&val).Positive()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt16Positive, ErrCause(verr))
}

func TestInt16Negative(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int16 = -5
	verr = Int16(&val).Negative()()
	assert.Nil(verr)

	verr = Int16(nil).Negative()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt16Negative, ErrCause(verr))

	val = 5
	verr = Int16(&val).Negative()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt16Negative, ErrCause(verr))
}

func TestInt16Zero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int16 = 0
	verr = Int16(&val).Zero()()
	assert.Nil(verr)

	verr = Int16(nil).Zero()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt16Zero, ErrCause(verr))

	val = 5
	verr = Int16(&val).Zero()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt16Zero, ErrCause(verr))
}

func TestInt16NotZero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int16 = 5
	verr = Int16(&val).NotZero()()
	assert.Nil(verr)

	verr = Int16(nil).NotZero()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt16NotZero, ErrCause(verr))

	val = 0
	verr = Int16(&val).NotZero()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt16NotZero, ErrCause(verr))
}
