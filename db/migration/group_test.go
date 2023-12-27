/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
