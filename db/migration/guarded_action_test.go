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
