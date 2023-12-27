/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package validate

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestUint32Min(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val uint32 = 10
	verr = Uint32(&val).Min(1)()
	assert.Nil(verr)

	val = 10
	verr = Uint32(&val).Min(10)()
	assert.Nil(verr)

	verr = Uint32(nil).Min(10)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrUint32Min, ErrCause(verr))

	val = 1
	verr = Uint32(&val).Min(10)()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrUint32Min, ErrCause(verr))

	val = 10
	verr = Uint32(&val).Min(10)()
	assert.Nil(verr)
}

func TestUint32Max(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val uint32 = 1
	verr = Uint32(&val).Max(10)()
	assert.Nil(verr)

	verr = Uint32(nil).Max(10)()
	assert.Nil(verr)

	val = 10
	verr = Uint32(&val).Max(10)()
	assert.Nil(verr)

	val = 11
	verr = Uint32(&val).Max(10)()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrUint32Max, ErrCause(verr))
}

func TestUint32Between(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val uint32 = 5
	verr = Uint32(&val).Between(1, 10)()
	assert.Nil(verr)

	verr = Uint32(nil).Between(5, 10)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrUint32Min, ErrCause(verr))

	val = 1
	verr = Uint32(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(1, ErrValue(verr))
	assert.Equal(ErrUint32Min, ErrCause(verr))

	val = 5
	verr = Uint32(&val).Between(5, 10)()
	assert.Nil(verr)

	val = 10
	verr = Uint32(&val).Between(5, 10)()
	assert.Nil(verr)

	val = 11
	verr = Uint32(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(11, ErrValue(verr))
	assert.Equal(ErrUint32Max, ErrCause(verr))
}

func TestUint32Zero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val uint32 = 0
	verr = Uint32(&val).Zero()()
	assert.Nil(verr)

	verr = Uint32(nil).Zero()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrUint32Zero, ErrCause(verr))

	val = 5
	verr = Uint32(&val).Zero()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrUint32Zero, ErrCause(verr))
}

func TestUint32NotZero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val uint32 = 5
	verr = Uint32(&val).NotZero()()
	assert.Nil(verr)

	verr = Uint32(nil).NotZero()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrUint32NotZero, ErrCause(verr))

	val = 0
	verr = Uint32(&val).NotZero()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrUint32NotZero, ErrCause(verr))
}
