package cron

import "time"

// JobConfig is a configuration set for a job.
type JobConfig struct {
	// Labels define extra metadata that can be used to filter jobs.
	Labels map[string]string `json:"labels" yaml:"labels"`
	// Timeout represents the abort threshold for the job.
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
	// ShutdownGracePeriod represents the time a job is given to clean itself up.
	ShutdownGracePeriod time.Duration `json:"shutdownGracePeriod" yaml:"shutdownGracePeriod"`
	// Serial indicates if job executions cannot overlap.
	Serial *bool `json:"serial" yaml:"serial"`
	// HistoryEnabled sets if we should save invocation history and restore it.
	HistoryEnabled *bool `json:"historyEnabled" yaml:"historyEnabled"`
	// HistoryMaxCount is the maximum number of history items to keep.
	HistoryMaxCount int `json:"historyMaxCount" yaml:"historyMaxCount"`
	// HistoryMaxAge is the maximum age of history items to keep.
	HistoryMaxAge time.Duration `json:"historyMaxAge" yaml:"historyMaxAge"`
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

// SerialOrDefault returns a value or a default.
func (jc JobConfig) SerialOrDefault() bool {
	if jc.Serial != nil {
		return *jc.Serial
	}
	return DefaultSerial
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
	return DefaultHistoryMaxCount
}

// HistoryMaxAgeOrDefault returns a value or a default.
func (jc JobConfig) HistoryMaxAgeOrDefault() time.Duration {
	if jc.HistoryMaxAge > 0 {
		return jc.HistoryMaxAge
	}
	return DefaultHistoryMaxAge
}
