package cron

import "time"

// JobSchedulerStats represent stats about a job scheduler.
type JobSchedulerStats struct {
	JobName        string        `json:"jobName"`
	SuccessRate    float64       `json:"successRate"`
	OutputBytes    int           `json:"outputBytes"`
	RunsTotal      int           `json:"runsTotal"`
	RunsSuccessful int           `json:"runsSuccessful"`
	RunsFailed     int           `json:"runsFailed"`
	RunsCancelled  int           `json:"runsCancelled"`
	RunsTimedOut   int           `json:"runsTimedOut"`
	ElapsedMax     time.Duration `json:"elapsedMax"`
	ElapsedMin     time.Duration `json:"elapsedMin"`
	Elapsed50th    time.Duration `json:"elapsed50th"`
	Elapsed95th    time.Duration `json:"elapsed95th"`
}
