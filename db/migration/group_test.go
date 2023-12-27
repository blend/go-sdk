/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package migration

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNewGroup(t *testing.T) {
	assert := assert.New(t)

	g := NewGroup(OptGroupActions(NewStep(Always(), ActionFunc(NoOp))))
	assert.Len(g.Actions, 1)
}
