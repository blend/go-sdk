/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestInt32Ptr(t *testing.T) {
	assert := assert.New(t)

	isNil := Int32Ptr(nil)
	var value int32 = 1
	hasValue := Int32Ptr(&value)
	var value2 int32 = 2
	hasValue2 := Int32Ptr(&value2)

	var setValue int32
	assert.Nil(SetInt32(&setValue, isNil, hasValue, hasValue2)(context.TODO()))
	assert.Equal(1, setValue)
}
