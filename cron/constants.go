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
	"time"
)

// Constats and defaults
const (
	DefaultTimeout			time.Duration	= 0
	DefaultHistoryRestoreTimeout			= 5 * time.Second
	DefaultShutdownGracePeriod	time.Duration	= 0
)

const (
	// DefaultDisabled is a default.
	DefaultDisabled	= false
	// DefaultShouldSkipLoggerListeners is a default.
	DefaultShouldSkipLoggerListeners	= false
	// DefaultShouldSkipLoggerOutput is a default.
	DefaultShouldSkipLoggerOutput	= false
)

const (
	// FlagBegin is an event flag.
	FlagBegin	= "cron.begin"
	// FlagComplete is an event flag.
	FlagComplete	= "cron.complete"
	// FlagSuccess is an event flag.
	FlagSuccess	= "cron.success"
	// FlagErrored is an event flag.
	FlagErrored	= "cron.errored"
	// FlagCanceled is an event flag.
	FlagCanceled	= "cron.canceled"
	// FlagBroken is an event flag.
	FlagBroken	= "cron.broken"
	// FlagFixed is an event flag.
	FlagFixed	= "cron.fixed"
	// FlagEnabled is an event flag.
	FlagEnabled	= "cron.enabled"
	// FlagDisabled is an event flag.
	FlagDisabled	= "cron.disabled"
)

// JobManagerState is a job manager status.
type JobManagerState string

// JobManagerState values.
const (
	JobManagerStateUnknown	JobManagerState	= "unknown"
	JobManagerStateRunning	JobManagerState	= "started"
	JobManagerStateStopped	JobManagerState	= "stopped"
)

// JobSchedulerState is a job manager status.
type JobSchedulerState string

// JobManagerState values.
const (
	JobSchedulerStateUnknown	JobSchedulerState	= "unknown"
	JobSchedulerStateRunning	JobSchedulerState	= "started"
	JobSchedulerStateStopped	JobSchedulerState	= "stopped"
)

// JobInvocationStatus is a job status.
type JobInvocationStatus string

// JobInvocationState values.
const (
	JobInvocationStatusIdle		JobInvocationStatus	= "idle"
	JobInvocationStatusRunning	JobInvocationStatus	= "running"
	JobInvocationStatusCanceled	JobInvocationStatus	= "canceled"
	JobInvocationStatusErrored	JobInvocationStatus	= "errored"
	JobInvocationStatusSuccess	JobInvocationStatus	= "success"
)
