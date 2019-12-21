package cron

import "time"

// JobSchedulerStatus is a status for a job scheduler.
type JobSchedulerStatus struct {
	Name        string            `json:"name"`
	State       JobSchedulerState `json:"state"`
	Labels      map[string]string `json:"labels"`
	Schedule    string            `json:"schedule"`
	Timeout     time.Duration     `json:"timeout"`
	Disabled    bool              `json:"disabled"`
	NextRuntime time.Time         `json:"nextRuntime"`
	Current     *JobInvocation    `json:"current"`
	Last        *JobInvocation    `json:"last"`
	Stats       JobSchedulerStats `json:"stats"`

	HistoryEnabled            bool          `json:"historyEnabled"`
	HistoryPersistenceEnabled bool          `json:"historyPersistenceEnabled"`
	HistoryMaxCount           int           `json:"historyMaxCount"`
	HistoryMaxAge             time.Duration `json:"historyMaxAge"`

	History []JobSchedulerStatusHistory `json:"history"`
}

// JobSchedulerStatusHistory is a state and elapsed pair used for graphs.
type JobSchedulerStatusHistory struct {
	State   JobInvocationState
	Elapsed time.Duration
}
