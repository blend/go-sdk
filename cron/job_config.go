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
	// HistoryPath is the path to write history to.
	HistoryPath string `json:"historyPath" yaml:"historyPath"`
	// HistoryMaxCount is the maximum number of history items to keep.
	HistoryMaxCount int `json:"historyMaxCount" yaml:"historyMaxCount"`
	// HistoryMaxAge is the maximum age of history items to keep.
	HistoryMaxAge int `json:"historyMaxAge" yaml:"historyMaxAge"`
}
