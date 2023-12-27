/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package db

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestError(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(Error(nil))

	var err error
	assert.Nil(Error(err))

	err = ex.New("this is only a test")
	assert.True(ex.Is(Error(err), ex.Class("this is only a test")))
}
