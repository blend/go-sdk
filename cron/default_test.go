/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
