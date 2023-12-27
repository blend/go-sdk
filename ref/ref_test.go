/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ref

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestRef(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(Ref("foo"))
	assert.NotEmpty(Refs("foo", "bar"))
	assert.NotNil(Ref(1))
	assert.NotEmpty(Refs(1, 2))
	assert.NotNil(Ref(time.Now()))
	assert.NotEmpty(Refs(time.Time{}, time.Time{}))

	assert.NotNil(String("foo"))
	assert.NotEmpty(Strings("foo", "bar"))
	assert.NotNil(Bool(true))
	assert.NotNil(Byte('b'))
	assert.NotNil(Rune('b'))
	assert.NotNil(Uint8(0))
	assert.NotNil(Uint16(0))
	assert.NotNil(Uint32(0))
	assert.NotNil(Uint64(0))
	assert.NotNil(Int8(0))
	assert.NotNil(Int16(0))
	assert.NotNil(Int32(0))
	assert.NotNil(Int64(0))
	assert.NotNil(Int(0))
	assert.NotNil(Float32(0))
	assert.NotNil(Float64(0))
	assert.NotNil(Time(time.Time{}))
	assert.NotNil(Duration(0))
}

func TestDeref(t *testing.T) {
	assert := assert.New(t)

	// strings
	populatedString := "Hello, world"
	emptyString := ""
	var nilStringPtr *string
	assert.Equal(populatedString, Deref(&populatedString))
	assert.Equal(emptyString, Deref(&emptyString))
	assert.Equal(emptyString, Deref(nilStringPtr))

	// ints
	populatedInt := 10
	zeroInt := 0
	var nilIntPtr *int
	assert.Equal(populatedInt, Deref(&populatedInt))
	assert.Equal(zeroInt, Deref(&zeroInt))
	assert.Equal(zeroInt, Deref(nilIntPtr))

	// times
	populatedTime := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	zeroTime := time.Time{}
	var nilTimePtr *time.Time
	assert.Equal(populatedTime, Deref(&populatedTime))
	assert.Equal(zeroTime, Deref(&zeroTime))
	assert.Equal(zeroTime, Deref(nilTimePtr))

	// slice of pointers
	assert.Equal(Derefs(&populatedInt, nil, &zeroInt), []int{populatedInt, zeroInt, zeroInt})
}
