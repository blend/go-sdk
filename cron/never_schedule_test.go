/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
