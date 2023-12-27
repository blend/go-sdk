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

func TestInt(t *testing.T) {
	assert := assert.New(t)

	intValue := Int(0)
	ptr, err := intValue.Int(context.TODO())
	assert.Nil(ptr)
	assert.Nil(err)

	intValue = Int(1234)
	ptr, err = intValue.Int(context.TODO())
	assert.Nil(err)
	assert.NotNil(ptr)
	assert.Equal(1234, *ptr)
}
