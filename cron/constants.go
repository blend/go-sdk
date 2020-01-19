package cron

import (
	"time"
)

// Constats and defaults
const (
	DefaultTimeout               time.Duration = 0
	DefaultHistoryRestoreTimeout               = 5 * time.Second
	DefaultShutdownGracePeriod   time.Duration = 0
)

const (
	// DefaultDisabled is a default.
	DefaultDisabled = false
	// DefaultShouldSkipLoggerListeners is a default.
	DefaultShouldSkipLoggerListeners = false
	// DefaultShouldSkipLoggerOutput is a default.
	DefaultShouldSkipLoggerOutput = false
	// DefaultHistoryEnabled is a default.
	DefaultHistoryEnabled = false
	// DefaultHistoryPersistenceEnabled is a default.
	DefaultHistoryPersistenceEnabled = false
	// DefaultHistoryMaxCount is the default number of history items to track.
	DefaultHistoryMaxCount = 10
	// DefaultHistoryMaxAge is the default maximum age of history items.
	DefaultHistoryMaxAge = 6 * time.Hour
)

const (
	// FlagBegin is an event flag.
	FlagBegin = "cron.begin"
	// FlagComplete is an event flag.
	FlagComplete = "cron.complete"
	// FlagSuccess is an event flag.
	FlagSuccess = "cron.success"
	// FlagErrored is an event flag.
	FlagErrored = "cron.errored"
	// FlagCancelled is an event flag.
	FlagCancelled = "cron.cancelled"
	// FlagBroken is an event flag.
	FlagBroken = "cron.broken"
	// FlagFixed is an event flag.
	FlagFixed = "cron.fixed"
	// FlagEnabled is an event flag.
	FlagEnabled = "cron.enabled"
	// FlagDisabled is an event flag.
	FlagDisabled = "cron.disabled"
)

// JobManagerState is a job manager status.
type JobManagerState string

// JobManagerState values.
const (
	JobManagerStateUnknown JobManagerState = "unknown"
	JobManagerStateRunning JobManagerState = "started"
	JobManagerStateStopped JobManagerState = "stopped"
)

// JobSchedulerState is a job manager status.
type JobSchedulerState string

// JobManagerState values.
const (
	JobSchedulerStateUnknown JobSchedulerState = "unknown"
	JobSchedulerStateRunning JobSchedulerState = "started"
	JobSchedulerStateStopped JobSchedulerState = "stopped"
)

// JobInvocationStatus is a job status.
type JobInvocationStatus string

// JobInvocationState values.
const (
	JobInvocationStatusIdle      JobInvocationStatus = "idle"
	JobInvocationStatusRunning   JobInvocationStatus = "running"
	JobInvocationStatusCancelled JobInvocationStatus = "cancelled"
	JobInvocationStatusErrored   JobInvocationStatus = "errored"
	JobInvocationStatusSuccess   JobInvocationStatus = "success"
)
