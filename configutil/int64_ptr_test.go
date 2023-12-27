/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestInt64Ptr(t *testing.T) {
	assert := assert.New(t)

	isNil := Int64Ptr(nil)
	value := int64(1)
	hasValue := Int64Ptr(&value)
	value2 := int64(2)
	hasValue2 := Int64Ptr(&value2)

	var setValue int64
	assert.Nil(SetInt64(&setValue, isNil, hasValue, hasValue2)(context.TODO()))
	assert.Equal(1, setValue)
}
