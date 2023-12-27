/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package cron

import (
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func Test_DelaySchedule(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	ts := time.Date(2019, 9, 8, 12, 11, 10, 9, time.UTC)
	ds := Delay(500*time.Millisecond, EverySecond())

	next := ds.Next(ts)
	its.Equal(ts.Add(500*time.Millisecond).Add(time.Second), next)
	its.Equal(1, ds.didRun)

	next = ds.Next(ts)
	its.Equal(ts.Add(time.Second), next)
	its.Equal(1, ds.didRun)

	// do this again to stress the `didRun` acas
	next = ds.Next(ts)
	its.Equal(ts.Add(time.Second), next)
	its.Equal(1, ds.didRun)
}

func Test_DelaySchedule_parse(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	ds := Delay(500*time.Millisecond, EverySecond())

	parsed, err := ParseSchedule(ds.String())
	its.Nil(err)
	its.Equal(fmt.Sprint(ds), fmt.Sprint(parsed))
}
