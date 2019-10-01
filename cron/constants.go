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
	// DefaultHistoryDisabled is a default.
	DefaultHistoryDisabled = false
	// DefaultHistoryPersistenceDisabled is a default.
	DefaultHistoryPersistenceDisabled = false
	// DefaultHistoryMaxCount is the default number of history items to track.
	DefaultHistoryMaxCount = 10
	// DefaultHistoryMaxAge is the default maximum age of history items.
	DefaultHistoryMaxAge = 6 * time.Hour
)

const (
	// FlagStarted is an event flag.
	FlagStarted = "cron.started"
	// FlagFailed is an event flag.
	FlagFailed = "cron.failed"
	// FlagCancelled is an event flag.
	FlagCancelled = "cron.cancelled"
	// FlagComplete is an event flag.
	FlagComplete = "cron.complete"
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
	JobManagerStatePaused  JobManagerState = "paused"
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

// JobInvocationState is a job status.
type JobInvocationState string

// JobInvocationState values.
const (
	JobInvocationStateRunning   JobInvocationState = "running"
	JobInvocationStateCancelled JobInvocationState = "cancelled"
	JobInvocationStateFailed    JobInvocationState = "failed"
	JobInvocationStateComplete  JobInvocationState = "complete"
)
