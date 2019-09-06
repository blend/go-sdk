package cron

import "time"

// Status is a status object
type Status struct {
	Status         JobManagerStatus            `json:"state"`
	Started        time.Time                   `json:"started"`
	Paused         time.Time                   `json:"paused"`
	Stopped        time.Time                   `json:"stopped"`
	JobLastStarted time.Time                   `json:"jobLastStarted"`
	Jobs           []*JobScheduler             `json:"jobs"`
	Running        map[string][]*JobInvocation `json:"running,omitempty"`
}
