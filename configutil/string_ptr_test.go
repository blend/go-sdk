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

func TestStringPtr(t *testing.T) {
	assert := assert.New(t)

	isNil := StringPtr(nil)
	emptyValue := ""
	isEmpty := StringPtr(&emptyValue)
	value := "foo"
	hasValue := StringPtr(&value)
	value2 := "bar"
	hasValue2 := StringPtr(&value2)

	var setValue string
	assert.Nil(SetString(&setValue, isNil, isEmpty, hasValue, hasValue2)(context.TODO()))
	assert.Equal("", setValue)

	assert.Nil(SetString(&setValue, isNil, hasValue, hasValue2)(context.TODO()))
	assert.Equal(value, setValue)
}
