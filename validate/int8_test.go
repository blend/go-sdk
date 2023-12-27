/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package validate

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestInt8Min(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int8 = 10
	verr = Int8(&val).Min(1)()
	assert.Nil(verr)

	val = 10
	verr = Int8(&val).Min(10)()
	assert.Nil(verr)

	verr = Int8(nil).Min(10)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt8Min, ErrCause(verr))

	val = 1
	verr = Int8(&val).Min(10)()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt8Min, ErrCause(verr))

	val = 10
	verr = Int8(&val).Min(10)()
	assert.Nil(verr)
}

func TestInt8Max(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int8 = 1
	verr = Int8(&val).Max(10)()
	assert.Nil(verr)

	verr = Int8(nil).Max(10)()
	assert.Nil(verr)

	val = 10
	verr = Int8(&val).Max(10)()
	assert.Nil(verr)

	val = 11
	verr = Int8(&val).Max(10)()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt8Max, ErrCause(verr))
}

func TestInt8Between(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int8 = 5
	verr = Int8(&val).Between(1, 10)()
	assert.Nil(verr)

	verr = Int8(nil).Between(5, 10)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt8Min, ErrCause(verr))

	val = 1
	verr = Int8(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(1, ErrValue(verr))
	assert.Equal(ErrInt8Min, ErrCause(verr))

	val = 5
	verr = Int8(&val).Between(5, 10)()
	assert.Nil(verr)

	val = 10
	verr = Int8(&val).Between(5, 10)()
	assert.Nil(verr)

	val = 11
	verr = Int8(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(11, ErrValue(verr))
	assert.Equal(ErrInt8Max, ErrCause(verr))
}

func TestInt8Positive(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int8 = 5
	verr = Int8(&val).Positive()()
	assert.Nil(verr)

	verr = Int8(nil).Positive()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt8Positive, ErrCause(verr))

	val = -5
	verr = Int8(&val).Positive()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt8Positive, ErrCause(verr))
}

func TestInt8Negative(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int8 = -5
	verr = Int8(&val).Negative()()
	assert.Nil(verr)

	verr = Int8(nil).Negative()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt8Negative, ErrCause(verr))

	val = 5
	verr = Int8(&val).Negative()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt8Negative, ErrCause(verr))
}

func TestInt8Zero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int8 = 0
	verr = Int8(&val).Zero()()
	assert.Nil(verr)

	verr = Int8(nil).Zero()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt8Zero, ErrCause(verr))

	val = 5
	verr = Int8(&val).Zero()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt8Zero, ErrCause(verr))
}

func TestInt8NotZero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int8 = 5
	verr = Int8(&val).NotZero()()
	assert.Nil(verr)

	verr = Int8(nil).NotZero()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt8NotZero, ErrCause(verr))

	val = 0
	verr = Int8(&val).NotZero()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt8NotZero, ErrCause(verr))
}
