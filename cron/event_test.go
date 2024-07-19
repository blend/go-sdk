/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cron

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNewEvent(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	e := NewEvent(FlagComplete, "test_task")
	its.Equal(FlagComplete, e.GetFlag())
	its.Equal("test_task", e.JobName)
}
