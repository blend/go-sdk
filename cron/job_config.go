package cron

import "time"

// JobConfig is a configuration set for a job.
type JobConfig struct {
	// Name sets the job name.
	Name string `json:"name" yaml:"name"`
	// Disabled determines if the job should be automatically scheduled or not.
	Disabled *bool `json:"disabled" yaml:"disabled"`
	// Description is an optional string to describe what the job does.
	Description string `json:"description" yaml:"description"`
	// Labels define extra metadata that can be used to filter jobs.
	Labels map[string]string `json:"labels" yaml:"labels"`
	// Timeout represents the abort threshold for the job.
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
	// ShutdownGracePeriod represents the time a job is given to clean itself up.
	ShutdownGracePeriod time.Duration `json:"shutdownGracePeriod" yaml:"shutdownGracePeriod"`
	// HistoryEnabled sets if we should save invocation history and restore it.
	HistoryEnabled *bool `json:"historyEnabled" yaml:"historyEnabled"`
	// HistoryPersistenceEnabled determines if we should call the history persister if one is provided.
	HistoryPersistenceEnabled *bool `json:"historyPersistenceEnabled" yaml:"historyPersistenceEnabled"`
	// HistoryMaxCount is the maximum number of history items to keep.
	HistoryMaxCount int `json:"historyMaxCount" yaml:"historyMaxCount"`
	// HistoryMaxAge is the maximum age of history items to keep.
	HistoryMaxAge time.Duration `json:"historyMaxAge" yaml:"historyMaxAge"`

	// ShouldSkipLoggerListeners skips triggering logger events if it is set to true.
	ShouldSkipLoggerListeners *bool `json:"shouldSkipLoggerListeners" yaml:"shouldSkipLoggerListeners"`
	// ShouldSkipLoggerOutput skips writing logger output if it is set to true.
	ShouldSkipLoggerOutput *bool `json:"shouldSkipLoggerOutput" yaml:"shouldSkipLoggerOutput"`
}

// DisabledOrDefault returns a value or a default.
func (jc JobConfig) DisabledOrDefault() bool {
	if jc.Disabled != nil {
		return *jc.Disabled
	}
	return DefaultDisabled
}

// TimeoutOrDefault returns a value or a default.
func (jc JobConfig) TimeoutOrDefault() time.Duration {
	if jc.Timeout > 0 {
		return jc.Timeout
	}
	return DefaultTimeout
}

// ShutdownGracePeriodOrDefault returns a value or a default.
func (jc JobConfig) ShutdownGracePeriodOrDefault() time.Duration {
	if jc.ShutdownGracePeriod > 0 {
		return jc.ShutdownGracePeriod
	}
	return DefaultShutdownGracePeriod
}

// HistoryEnabledOrDefault returns a value or a default.
func (jc JobConfig) HistoryEnabledOrDefault() bool {
	if jc.HistoryEnabled != nil {
		return *jc.HistoryEnabled
	}
	return DefaultHistoryEnabled
}

// HistoryMaxCountOrDefault returns a value or a default.
func (jc JobConfig) HistoryMaxCountOrDefault() int {
	if jc.HistoryMaxCount > 0 {
		return jc.HistoryMaxCount
	}
	return 0
}

// HistoryMaxAgeOrDefault returns a value or a default.
func (jc JobConfig) HistoryMaxAgeOrDefault() time.Duration {
	if jc.HistoryMaxAge > 0 {
		return jc.HistoryMaxAge
	}
	return 0
}

// HistoryPersistenceEnabledOrDefault returns a value or a default.
func (jc JobConfig) HistoryPersistenceEnabledOrDefault() bool {
	if jc.HistoryPersistenceEnabled != nil {
		return *jc.HistoryPersistenceEnabled
	}
	return DefaultHistoryPersistenceEnabled
}

// ShouldSkipLoggerListenersOrDefault returns a value or a default.
func (jc JobConfig) ShouldSkipLoggerListenersOrDefault() bool {
	if jc.ShouldSkipLoggerListeners != nil {
		return *jc.ShouldSkipLoggerListeners
	}
	return DefaultShouldSkipLoggerListeners
}

// ShouldSkipLoggerOutputOrDefault returns a value or a default.
func (jc JobConfig) ShouldSkipLoggerOutputOrDefault() bool {
	if jc.ShouldSkipLoggerOutput != nil {
		return *jc.ShouldSkipLoggerOutput
	}
	return DefaultShouldSkipLoggerOutput
}
