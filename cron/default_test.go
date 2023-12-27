/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package cron

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestDefault(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(Default())

	SetDefault(nil)
	assert.Nil(_default)
}
