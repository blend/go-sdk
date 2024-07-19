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

func TestInt32(t *testing.T) {
	assert := assert.New(t)

	intValue := Int32(0)
	ptr, err := intValue.Int32(context.TODO())
	assert.Nil(ptr)
	assert.Nil(err)

	intValue = Int32(1234)
	ptr, err = intValue.Int32(context.TODO())
	assert.Nil(err)
	assert.NotNil(ptr)
	assert.Equal(1234, *ptr)
}
