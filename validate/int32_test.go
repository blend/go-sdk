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

func TestInt32Min(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int32 = 10
	verr = Int32(&val).Min(1)()
	assert.Nil(verr)

	val = 10
	verr = Int32(&val).Min(10)()
	assert.Nil(verr)

	verr = Int32(nil).Min(10)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt32Min, ErrCause(verr))

	val = 1
	verr = Int32(&val).Min(10)()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt32Min, ErrCause(verr))

	val = 10
	verr = Int32(&val).Min(10)()
	assert.Nil(verr)
}

func TestInt32Max(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int32 = 1
	verr = Int32(&val).Max(10)()
	assert.Nil(verr)

	verr = Int32(nil).Max(10)()
	assert.Nil(verr)

	val = 10
	verr = Int32(&val).Max(10)()
	assert.Nil(verr)

	val = 11
	verr = Int32(&val).Max(10)()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt32Max, ErrCause(verr))
}

func TestInt32Between(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int32 = 5
	verr = Int32(&val).Between(1, 10)()
	assert.Nil(verr)

	verr = Int32(nil).Between(5, 10)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt32Min, ErrCause(verr))

	val = 1
	verr = Int32(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(1, ErrValue(verr))
	assert.Equal(ErrInt32Min, ErrCause(verr))

	val = 5
	verr = Int32(&val).Between(5, 10)()
	assert.Nil(verr)

	val = 10
	verr = Int32(&val).Between(5, 10)()
	assert.Nil(verr)

	val = 11
	verr = Int32(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(11, ErrValue(verr))
	assert.Equal(ErrInt32Max, ErrCause(verr))
}

func TestInt32Positive(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int32 = 5
	verr = Int32(&val).Positive()()
	assert.Nil(verr)

	verr = Int32(nil).Positive()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt32Positive, ErrCause(verr))

	val = -5
	verr = Int32(&val).Positive()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt32Positive, ErrCause(verr))
}

func TestInt32Negative(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int32 = -5
	verr = Int32(&val).Negative()()
	assert.Nil(verr)

	verr = Int32(nil).Negative()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt32Negative, ErrCause(verr))

	val = 5
	verr = Int32(&val).Negative()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt32Negative, ErrCause(verr))
}

func TestInt32Zero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int32 = 0
	verr = Int32(&val).Zero()()
	assert.Nil(verr)

	verr = Int32(nil).Zero()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt32Zero, ErrCause(verr))

	val = 5
	verr = Int32(&val).Zero()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt32Zero, ErrCause(verr))
}

func TestInt32NotZero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int32 = 5
	verr = Int32(&val).NotZero()()
	assert.Nil(verr)

	verr = Int32(nil).NotZero()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt32NotZero, ErrCause(verr))

	val = 0
	verr = Int32(&val).NotZero()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt32NotZero, ErrCause(verr))
}
