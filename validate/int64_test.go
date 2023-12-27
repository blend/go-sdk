/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package validate

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestInt64Min(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int64 = 10
	verr = Int64(&val).Min(1)()
	assert.Nil(verr)

	val = 10
	verr = Int64(&val).Min(10)()
	assert.Nil(verr)

	verr = Int64(nil).Min(10)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt64Min, ErrCause(verr))

	val = 1
	verr = Int64(&val).Min(10)()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt64Min, ErrCause(verr))

	val = 10
	verr = Int64(&val).Min(10)()
	assert.Nil(verr)
}

func TestInt64Max(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int64 = 1
	verr = Int64(&val).Max(10)()
	assert.Nil(verr)

	verr = Int64(nil).Max(10)()
	assert.Nil(verr)

	val = 10
	verr = Int64(&val).Max(10)()
	assert.Nil(verr)

	val = 11
	verr = Int64(&val).Max(10)()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt64Max, ErrCause(verr))
}

func TestInt64Between(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int64 = 5
	verr = Int64(&val).Between(1, 10)()
	assert.Nil(verr)

	verr = Int64(nil).Between(5, 10)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt64Min, ErrCause(verr))

	val = 1
	verr = Int64(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(1, ErrValue(verr))
	assert.Equal(ErrInt64Min, ErrCause(verr))

	val = 5
	verr = Int64(&val).Between(5, 10)()
	assert.Nil(verr)

	val = 10
	verr = Int64(&val).Between(5, 10)()
	assert.Nil(verr)

	val = 11
	verr = Int64(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(11, ErrValue(verr))
	assert.Equal(ErrInt64Max, ErrCause(verr))
}

func TestInt64Positive(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int64 = 5
	verr = Int64(&val).Positive()()
	assert.Nil(verr)

	verr = Int64(nil).Positive()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt64Positive, ErrCause(verr))

	val = -5
	verr = Int64(&val).Positive()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt64Positive, ErrCause(verr))
}

func TestInt64Negative(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int64 = -5
	verr = Int64(&val).Negative()()
	assert.Nil(verr)

	verr = Int64(nil).Negative()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt64Negative, ErrCause(verr))

	val = 5
	verr = Int64(&val).Negative()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt64Negative, ErrCause(verr))
}

func TestInt64Zero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int64 = 0
	verr = Int64(&val).Zero()()
	assert.Nil(verr)

	verr = Int64(nil).Zero()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt64Zero, ErrCause(verr))

	val = 5
	verr = Int64(&val).Zero()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt64Zero, ErrCause(verr))
}

func TestInt64NotZero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val int64 = 5
	verr = Int64(&val).NotZero()()
	assert.Nil(verr)

	verr = Int64(nil).NotZero()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrInt64NotZero, ErrCause(verr))

	val = 0
	verr = Int64(&val).NotZero()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrInt64NotZero, ErrCause(verr))
}
