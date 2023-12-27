/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package slant

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestPrint(t *testing.T) {
	assert := assert.New(t)

	output, err := PrintString("WARDEN")
	assert.Nil(err)
	assert.NotEmpty(output)
}
