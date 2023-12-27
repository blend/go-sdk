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

func TestInt646464646464(t *testing.T) {
	assert := assert.New(t)

	intValue := Int64(0)
	ptr, err := intValue.Int64(context.TODO())
	assert.Nil(ptr)
	assert.Nil(err)

	intValue = Int64(1234)
	ptr, err = intValue.Int64(context.TODO())
	assert.Nil(err)
	assert.NotNil(ptr)
	assert.Equal(1234, *ptr)
}
