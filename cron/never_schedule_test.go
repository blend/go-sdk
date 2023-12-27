/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package cron

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func Test_NeverSchedule(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	its.Equal(time.Time{}, new(NeverSchedule).Next(time.Now().UTC()))
	its.Equal(StringScheduleNever, new(NeverSchedule).String())
}
