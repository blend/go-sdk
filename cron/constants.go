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
	// DefaultEnabled is a default.
	DefaultEnabled = true
	// DefaultSerial is a default.
	DefaultSerial = false
	// DefaultShouldSkipLoggerListeners is a default.
	DefaultShouldSkipLoggerListeners = false
	// DefaultShouldSkipLoggerOutput is a default.
	DefaultShouldSkipLoggerOutput = false
	// DefaultHistoryEnabled is a default.
	DefaultHistoryEnabled = true
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

// JobManagerStatus is a job manager status.
type JobManagerStatus string

const (
	// JobManagerStatusRunning is a job manager status.
	JobManagerStatusRunning JobManagerStatus = "started"
	// JobManagerStatusPaused is a job manager status.
	JobManagerStatusPaused JobManagerStatus = "paused"
	// JobManagerStatusResuming is a job manager status.
	JobManagerStatusResuming JobManagerStatus = "resuming"
	// JobManagerStatusStopped is a job manager status.
	JobManagerStatusStopped JobManagerStatus = "stopped"
)

// JobStatus is a job status.
type JobStatus string

// JobStatus values.
const (
	JobStatusRunning   JobStatus = "running"
	JobStatusCancelled JobStatus = "cancelled"
	JobStatusFailed    JobStatus = "failed"
	JobStatusComplete  JobStatus = "complete"
)
