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

func TestLazyInt64(t *testing.T) {
	its := assert.New(t)

	isNil := LazyInt64(nil)
	var value int64 = 0
	hasValue := LazyInt64(&value)
	var value2 int64 = 2
	hasValue2 := LazyInt64(&value2)

	var setValue int64
	its.Nil(SetInt64(&setValue, isNil, hasValue, hasValue2)(context.TODO()))
	its.Equal(2, setValue)
}
