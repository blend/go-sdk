/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package validate

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestUint16Min(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val uint16 = 10
	verr = Uint16(&val).Min(1)()
	assert.Nil(verr)

	val = 10
	verr = Uint16(&val).Min(10)()
	assert.Nil(verr)

	verr = Uint16(nil).Min(10)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrUint16Min, ErrCause(verr))

	val = 1
	verr = Uint16(&val).Min(10)()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrUint16Min, ErrCause(verr))

	val = 10
	verr = Uint16(&val).Min(10)()
	assert.Nil(verr)
}

func TestUint16Max(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val uint16 = 1
	verr = Uint16(&val).Max(10)()
	assert.Nil(verr)

	verr = Uint16(nil).Max(10)()
	assert.Nil(verr)

	val = 10
	verr = Uint16(&val).Max(10)()
	assert.Nil(verr)

	val = 11
	verr = Uint16(&val).Max(10)()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrUint16Max, ErrCause(verr))
}

func TestUint16Between(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val uint16 = 5
	verr = Uint16(&val).Between(1, 10)()
	assert.Nil(verr)

	verr = Uint16(nil).Between(5, 10)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrUint16Min, ErrCause(verr))

	val = 1
	verr = Uint16(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(1, ErrValue(verr))
	assert.Equal(ErrUint16Min, ErrCause(verr))

	val = 5
	verr = Uint16(&val).Between(5, 10)()
	assert.Nil(verr)

	val = 10
	verr = Uint16(&val).Between(5, 10)()
	assert.Nil(verr)

	val = 11
	verr = Uint16(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(11, ErrValue(verr))
	assert.Equal(ErrUint16Max, ErrCause(verr))
}

func TestUint16Zero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val uint16 = 0
	verr = Uint16(&val).Zero()()
	assert.Nil(verr)

	verr = Uint16(nil).Zero()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrUint16Zero, ErrCause(verr))

	val = 5
	verr = Uint16(&val).Zero()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrUint16Zero, ErrCause(verr))
}

func TestUint16NotZero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val uint16 = 5
	verr = Uint16(&val).NotZero()()
	assert.Nil(verr)

	verr = Uint16(nil).NotZero()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrUint16NotZero, ErrCause(verr))

	val = 0
	verr = Uint16(&val).NotZero()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrUint16NotZero, ErrCause(verr))
}
