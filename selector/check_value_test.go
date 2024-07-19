/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package selector

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestCheckValue(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(CheckValue(""), "should not error on empty values")
	assert.Nil(CheckValue("foo"))
	assert.Nil(CheckValue("bar_baz"))
	assert.NotNil(CheckValue("_bar_baz"))
	assert.NotNil(CheckValue("bar_baz_"))
	assert.NotNil(CheckValue("_bar_baz_"))
}
