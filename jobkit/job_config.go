package jobkit

import (
	"path/filepath"
	"time"

	"github.com/blend/go-sdk/email"
	"github.com/blend/go-sdk/stringutil"
)

// JobConfig is something you can use to give your jobs some knobs to turn
// from configuration.
// You can use this job config by embedding it into your larger job config struct.
type JobConfig struct {
	// Name is the name of the job.
	Name string `json:"name" yaml:"name"`
	// Description is a description of the job.
	Description string `json:"description" yaml:"description"`
	// Labels define extra metadata that can be used to filter jobs.
	Labels map[string]string `json:"labels" yaml:"labels"`
	// Schedule returns the job schedule.
	Schedule string `json:"schedule" yaml:"schedule"`
	// Timeout represents the abort threshold for the job.
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
	// ShutdownGracePeriod represents the time a job is given to clean itself up.
	ShutdownGracePeriod time.Duration `json:"shutdownGracePeriod" yaml:"shutdownGracePeriod"`
	// Serial indicates if job executions cannot overlap.
	Serial *bool `json:"serial" yaml:"serial"`
	// HistoryEnabled sets if we should save invocation history and restore it.
	HistoryEnabled *bool `json:"historyEnabled" yaml:"historyEnabled"`
	// HistoryPath is the path to write history to.
	HistoryPath string `json:"historyPath" yaml:"historyPath"`
	// HistoryMaxCount is the maximum number of history items to keep.
	HistoryMaxCount int `json:"historyMaxCount" yaml:"historyMaxCount"`
	// HistoryMaxAge is the maximum age of history items to keep.
	HistoryMaxAge int `json:"historyMaxAge" yaml:"historyMaxAge"`

	// NotifyOnStart governs if we should send notifications job start.
	NotifyOnStart *bool `json:"notifyOnStart" yaml:"notifyOnStart"`
	// NotifyOnSuccess governs if we should send notifications on any success.
	NotifyOnSuccess *bool `json:"notifyOnSuccess" yaml:"notifyOnSuccess"`
	// NotifyOnFailure governs if we should send notifications on any failure.
	NotifyOnFailure *bool `json:"notifyOnFailure" yaml:"notifyOnFailure"`
	// NotifyOnCancellation governs if we should send notifications on cancellation.
	NotifyOnCancellation *bool `json:"notifyOnCancellation" yaml:"notifyOnCancellation"`
	// NotifyOnBroken governs if we should send notifications on a success => failure transition.
	NotifyOnBroken *bool `json:"notifyOnBroken" yaml:"notifyOnBroken"`
	// NotifyOnFixed governs if we should send notifications on a failure => success transition.
	NotifyOnFixed *bool `json:"notifyOnFixed" yaml:"notifyOnFixed"`
	// NotifyOnEnabled governs if we should send notifications when a job is enabled.
	NotifyOnEnabled *bool `json:"notifyOnEnabled" yaml:"notifyOnEnabled"`
	// NotifyOnDisabled governs if we should send notifications when a job is disabled.
	NotifyOnDisabled *bool `json:"notifyOnDisabled" yaml:"notifyOnDisabled"`

	// EmailDefaults are the message defaults for email notifications.
	EmailDefaults email.Message `json:"emailDefaults" yaml:"emailDefaults"`
}

// ScheduleOrDefault returns the schedule or a default (every 5 minutes).
func (jc JobConfig) ScheduleOrDefault() string {
	if jc.Schedule != "" {
		return jc.Schedule
	}
	return "* */5 * * * * *"
}

// SerialOrDefault returns a value or a default.
func (jc JobConfig) SerialOrDefault() bool {
	if jc.Serial != nil {
		return *jc.Serial
	}
	return true
}

// HistoryEnabledOrDefault returns a value or a default.
func (jc JobConfig) HistoryEnabledOrDefault() bool {
	if jc.HistoryEnabled != nil {
		return *jc.HistoryEnabled
	}
	return true
}

// HistoryPathOrDefault returns a value or a default.
func (jc JobConfig) HistoryPathOrDefault() string {
	if jc.HistoryPath != "" {
		return jc.HistoryPath
	}
	return filepath.Join("_history", stringutil.Slugify(jc.Name)+".json")
}

// NotifyOnStartOrDefault returns a value or a default.
func (jc JobConfig) NotifyOnStartOrDefault() bool {
	if jc.NotifyOnStart != nil {
		return *jc.NotifyOnStart
	}
	return false
}

// NotifyOnSuccessOrDefault returns a value or a default.
func (jc JobConfig) NotifyOnSuccessOrDefault() bool {
	if jc.NotifyOnSuccess != nil {
		return *jc.NotifyOnSuccess
	}
	return false
}

// NotifyOnFailureOrDefault returns a value or a default.
func (jc JobConfig) NotifyOnFailureOrDefault() bool {
	if jc.NotifyOnFailure != nil {
		return *jc.NotifyOnFailure
	}
	return true
}

// NotifyOnCancellationOrDefault returns a value or a default.
func (jc JobConfig) NotifyOnCancellationOrDefault() bool {
	if jc.NotifyOnCancellation != nil {
		return *jc.NotifyOnCancellation
	}
	return true
}

// NotifyOnBrokenOrDefault returns a value or a default.
func (jc JobConfig) NotifyOnBrokenOrDefault() bool {
	if jc.NotifyOnBroken != nil {
		return *jc.NotifyOnBroken
	}
	return true
}

// NotifyOnFixedOrDefault returns a value or a default.
func (jc JobConfig) NotifyOnFixedOrDefault() bool {
	if jc.NotifyOnFixed != nil {
		return *jc.NotifyOnFixed
	}
	return true
}

// NotifyOnEnabledOrDefault returns a value or a default.
func (jc JobConfig) NotifyOnEnabledOrDefault() bool {
	if jc.NotifyOnEnabled != nil {
		return *jc.NotifyOnEnabled
	}
	return false
}

// NotifyOnDisabledOrDefault returns a value or a default.
func (jc JobConfig) NotifyOnDisabledOrDefault() bool {
	if jc.NotifyOnDisabled != nil {
		return *jc.NotifyOnDisabled
	}
	return false
}
