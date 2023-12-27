/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
