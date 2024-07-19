/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package validate

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestReturnAll(t *testing.T) {
	assert := assert.New(t)

	res := ReturnAll(none, some(fmt.Errorf("one")), some(fmt.Errorf("two")), none)
	assert.Len(res, 2)
}
