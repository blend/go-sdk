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

func TestLazybool(t *testing.T) {
	its := assert.New(t)

	isNil := LazyBool(nil)
	var value bool = false
	hasValue := LazyBool(&value)
	var value2 bool = true
	hasValue2 := LazyBool(&value2)

	var setValue bool
	its.Nil(SetBool(&setValue, isNil, hasValue, hasValue2)(context.TODO()))
	its.Equal(true, setValue)
}
