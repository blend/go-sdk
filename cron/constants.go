package cron

import (
	"time"
)

const (
	// EnvVarHeartbeatInterval is an environment variable name.
	EnvVarHeartbeatInterval = "CRON_HEARTBEAT_INTERVAL"
)

// Constats and defaults
const (
	DefaultHistoryMaxCount                     = 10
	DefaultHistoryMaxAge                       = 6 * time.Hour
	DefaultHistoryRestoreTimeout               = 5 * time.Second
	DefaultShutdownGracePeriod   time.Duration = 0
	DefaultHeartbeatInterval                   = 50 * time.Millisecond
)

const (
	// DefaultEnabled is a default.
	DefaultEnabled = true
	// DefaultSerial is a default.
	DefaultSerial = false
	// DefaultShouldWriteOutput is a default.
	DefaultShouldWriteOutput = true
	// DefaultShouldTriggerListeners is a default.
	DefaultShouldTriggerListeners = true
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

// State is a job state.
type State string

const (
	//StateRunning is the running state.
	StateRunning State = "running"
	// StateEnabled is the enabled state.
	StateEnabled State = "enabled"
	// StateDisabled is the disabled state.
	StateDisabled State = "disabled"
)

// JobStatus is a job status.
type JobStatus string

// Status values.
const (
	JobStatusRunning   JobStatus = "running"
	JobStatusCancelled JobStatus = "cancelled"
	JobStatusFailed    JobStatus = "failed"
	JobStatusComplete  JobStatus = "complete"
)
