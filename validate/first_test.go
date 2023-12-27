/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package validate

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestFirst(t *testing.T) {
	assert := assert.New(t)

	res := First(none, some(fmt.Errorf("one")), some(fmt.Errorf("two")), none)()
	assert.Equal(fmt.Errorf("one"), res)
}
