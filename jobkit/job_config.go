package jobkit

import (
	"time"
)

// JobConfig is something you can use to give your jobs some knobs to turn
// from configuration.
// You can use this job config by embedding it into your larger job config struct.
type JobConfig struct {
	// Name is the name of the job.
	Name string `json:"name" yaml:"name"`
	// Description is a description of the job.
	Description string `json:"description" yaml:"description"`
	// Schedule returns the job schedule.
	Schedule string `json:"schedule" yaml:"schedule"`
	// Timeout represents the abort threshold for the job.
	Timeout time.Duration `json:"timeout" yaml:"timeout"`

	// NotifyOnStart governs if we should send notifications job start.
	NotifyOnStart *bool `json:"notifyOnStart" yaml:"notifyOnStart"`
	// NotifyOnSuccess governs if we should send notifications on any success.
	NotifyOnSuccess *bool `json:"notifyOnSuccess" yaml:"notifyOnSuccess"`
	// NotifyOnFailure governs if we should send notifications on any failure.
	NotifyOnFailure *bool `json:"notifyOnFailure" yaml:"notifyOnFailure"`
	// NotifyOnBroken governs if we should send notifications on a success => failure transition.
	NotifyOnBroken *bool `json:"notifyOnBroken" yaml:"notifyOnBroken"`
	// NotifyOnFixed governs if we should send notifications on a failure => success transition.
	NotifyOnFixed *bool `json:"notifyOnFixed" yaml:"notifyOnFixed"`
	// NotifyOnEnabled governs if we should send notifications when a job is enabled.
	NotifyOnEnabled *bool `json:"notifyOnEnabled" yaml:"notifyOnEnabled"`
	// NotifyOnDisabled governs if we should send notifications when a job is disabled.
	NotifyOnDisabled *bool `json:"notifyOnDisabled" yaml:"notifyOnDisabled"`
}

// ScheduleOrDefault returns the schedule or a default (every 5 minutes).
func (jc JobConfig) ScheduleOrDefault() string {
	if jc.Schedule != "" {
		return jc.Schedule
	}
	return "* */5 * * * * *"
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
