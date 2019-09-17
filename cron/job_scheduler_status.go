package cron

import "time"

// JobSchedulerStatus is a status for a job scheduler.
type JobSchedulerStatus struct {
	Name        string            `json:"name"`
	State       JobSchedulerState `json:"state"`
	Labels      map[string]string `json:"labels"`
	Schedule    string            `json:"schedule"`
	Disabled    bool              `json:"disabled"`
	NextRuntime time.Time         `json:"nextRuntime"`
	Current     *JobInvocation    `json:"current"`
	Last        *JobInvocation    `json:"last"`

	HistoryDisabled            bool          `json:"historyDisabled"`
	HistoryPersistenceDisabled bool          `json:"HistoryPersistenceDisabled"`
	HistoryMaxCount            int           `json:"historyMaxCount"`
	HistoryMaxAge              time.Duration `json:"historyMaxAge"`
}
