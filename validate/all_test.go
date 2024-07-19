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

func TestAll(t *testing.T) {
	assert := assert.New(t)

	res := All(none, some(fmt.Errorf("one")), some(fmt.Errorf("two")), none)()
	assert.Len(res, 2)
}

func TestAllNested(t *testing.T) {
	assert := assert.New(t)

	res := All(
		none,
		some(fmt.Errorf("one")),
		All(
			some(fmt.Errorf("two")),
			some(fmt.Errorf("three")),
		),
		none,
	)()
	assert.Len(res, 3)
}
