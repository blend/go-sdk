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

func TestLazyString(t *testing.T) {
	assert := assert.New(t)

	isNil := LazyString(nil)
	emptyValue := ""
	isEmpty := LazyString(&emptyValue)
	value := "foo"
	hasValue := LazyString(&value)
	value2 := "bar"
	hasValue2 := LazyString(&value2)

	var setValue string
	assert.Nil(SetString(&setValue, isNil, isEmpty, hasValue, hasValue2)(context.TODO()))
	assert.Equal(value, setValue)
}
