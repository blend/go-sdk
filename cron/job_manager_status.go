package cron

import "time"

// JobManagerStatus represents the status of a job manager.
type JobManagerStatus struct {
	State          JobManagerState          `json:"state"`
	Started        time.Time                `json:"started"`
	Paused         time.Time                `json:"paused"`
	Stopped        time.Time                `json:"stopped"`
	JobLastStarted time.Time                `json:"jobLastStarted"`
	Jobs           []JobSchedulerStatus     `json:"jobs"`
	Running        map[string]JobInvocation `json:"running"`
}
