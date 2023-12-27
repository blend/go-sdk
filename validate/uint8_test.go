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

func TestUint8Min(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val uint8 = 10
	verr = Uint8(&val).Min(1)()
	assert.Nil(verr)

	val = 10
	verr = Uint8(&val).Min(10)()
	assert.Nil(verr)

	verr = Uint8(nil).Min(10)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrUint8Min, ErrCause(verr))

	val = 1
	verr = Uint8(&val).Min(10)()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrUint8Min, ErrCause(verr))

	val = 10
	verr = Uint8(&val).Min(10)()
	assert.Nil(verr)
}

func TestUint8Max(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val uint8 = 1
	verr = Uint8(&val).Max(10)()
	assert.Nil(verr)

	verr = Uint8(nil).Max(10)()
	assert.Nil(verr)

	val = 10
	verr = Uint8(&val).Max(10)()
	assert.Nil(verr)

	val = 11
	verr = Uint8(&val).Max(10)()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrUint8Max, ErrCause(verr))
}

func TestUint8Between(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val uint8 = 5
	verr = Uint8(&val).Between(1, 10)()
	assert.Nil(verr)

	verr = Uint8(nil).Between(5, 10)()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrUint8Min, ErrCause(verr))

	val = 1
	verr = Uint8(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(1, ErrValue(verr))
	assert.Equal(ErrUint8Min, ErrCause(verr))

	val = 5
	verr = Uint8(&val).Between(5, 10)()
	assert.Nil(verr)

	val = 10
	verr = Uint8(&val).Between(5, 10)()
	assert.Nil(verr)

	val = 11
	verr = Uint8(&val).Between(5, 10)()
	assert.NotNil(verr)
	assert.Equal(11, ErrValue(verr))
	assert.Equal(ErrUint8Max, ErrCause(verr))
}

func TestUint8Zero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val uint8 = 0
	verr = Uint8(&val).Zero()()
	assert.Nil(verr)

	verr = Uint8(nil).Zero()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrUint8Zero, ErrCause(verr))

	val = 5
	verr = Uint8(&val).Zero()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrUint8Zero, ErrCause(verr))
}

func TestUint8NotZero(t *testing.T) {
	assert := assert.New(t)

	var verr error
	var val uint8 = 5
	verr = Uint8(&val).NotZero()()
	assert.Nil(verr)

	verr = Uint8(nil).NotZero()()
	assert.NotNil(verr)
	assert.Nil(ErrValue(verr))
	assert.Equal(ErrUint8NotZero, ErrCause(verr))

	val = 0
	verr = Uint8(&val).NotZero()()
	assert.NotNil(verr)
	assert.NotNil(ErrValue(verr))
	assert.Equal(ErrUint8NotZero, ErrCause(verr))
}
