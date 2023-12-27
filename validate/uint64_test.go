/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package validate

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestUint64Min(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val uint64 = 10
	verr = Uint64(&val).Min(1)()
	assert.Nil(verr)

	val = 10
	verr = Uint64(&val).Min(10)()
	assert.Nil(verr)

	verr = Uint64(nil).Min(10)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrUint64Min, ErrCause(verr))

	val = 1
	verr = Uint64(&val).Min(10)()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrUint64Min, ErrCause(verr))

	val = 10
	verr = Uint64(&val).Min(10)()
	assert.Nil(verr)
}

func TestUint64Max(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val uint64 = 1
	verr = Uint64(&val).Max(10)()
	assert.Nil(verr)

	verr = Uint64(nil).Max(10)()
	assert.Nil(verr)

	val = 10
	verr = Uint64(&val).Max(10)()
	assert.Nil(verr)

	val = 11
	verr = Uint64(&val).Max(10)()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrUint64Max, ErrCause(verr))
}

func TestUint64Between(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val uint64 = 5
	verr = Uint64(&val).Between(1, 10)()
	assert.Nil(verr)

	verr = Uint64(nil).Between(5, 10)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrUint64Min, ErrCause(verr))

	val = 1
	verr = Uint64(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(1, ErrValue(verr))
	assert.Equal(ErrUint64Min, ErrCause(verr))

	val = 5
	verr = Uint64(&val).Between(5, 10)()
	assert.Nil(verr)

	val = 10
	verr = Uint64(&val).Between(5, 10)()
	assert.Nil(verr)

	val = 11
	verr = Uint64(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(11, ErrValue(verr))
	assert.Equal(ErrUint64Max, ErrCause(verr))
}

func TestUint64Zero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val uint64 = 0
	verr = Uint64(&val).Zero()()
	assert.Nil(verr)

	verr = Uint64(nil).Zero()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrUint64Zero, ErrCause(verr))

	val = 5
	verr = Uint64(&val).Zero()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrUint64Zero, ErrCause(verr))
}

func TestUint64NotZero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val uint64 = 5
	verr = Uint64(&val).NotZero()()
	assert.Nil(verr)

	verr = Uint64(nil).NotZero()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrUint64NotZero, ErrCause(verr))

	val = 0
	verr = Uint64(&val).NotZero()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrUint64NotZero, ErrCause(verr))
}
