/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package migration

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestStep(t *testing.T) {
	assert := assert.New(t)

	step := NewStep(Always(), ActionFunc(NoOp))
	assert.NotNil(step.Guard)
	assert.NotNil(step.Body)
}
