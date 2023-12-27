/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package cron

import (
	"fmt"
	"time"
)

var (
	_	Schedule	= (*NeverSchedule)(nil)
	_	fmt.Stringer	= (*NeverSchedule)(nil)
)

// Never returns a never schedule.
func Never() NeverSchedule	{ return NeverSchedule{} }

// NeverSchedule is a schedule that never runs.
type NeverSchedule struct{}

// Next implements Schedule
func (ns NeverSchedule) Next(_ time.Time) time.Time	{ return time.Time{} }

// String implements fmt.Stringer.
func (ns NeverSchedule) String() string	{ return StringScheduleNever }
