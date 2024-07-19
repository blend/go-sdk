/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestIntPtr(t *testing.T) {
	assert := assert.New(t)

	isNil := IntPtr(nil)
	value := 1
	hasValue := IntPtr(&value)
	value2 := 2
	hasValue2 := IntPtr(&value2)

	var setValue int
	assert.Nil(SetInt(&setValue, isNil, hasValue, hasValue2)(context.TODO()))
	assert.Equal(1, setValue)
}

func TestIntPtr_Zero(t *testing.T) {
	assert := assert.New(t)

	isNil := IntPtr(nil)

	zero := 0
	zeroValue := IntPtr(&zero)

	value := 1
	hasValue := IntPtr(&value)
	value2 := 2
	hasValue2 := IntPtr(&value2)

	setValue := 3
	assert.Nil(SetInt(&setValue, isNil, zeroValue, hasValue, hasValue2)(context.TODO()))
	assert.Zero(setValue)
}
