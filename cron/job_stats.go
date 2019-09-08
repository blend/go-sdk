package cron

import "time"

// JobStats represent stats about a job based on history.
type JobStats struct {
	SuccessRate    float64       `json:"successRate"`
	OutputBytes    int           `json:"outputBytes"`
	RunsTotal      int           `json:"runsTotal"`
	RunsSuccessful int           `json:"runsSuccessful"`
	RunsCancelled  int           `json:"runsCancelled"`
	RunsTimedOut   int           `json:"runsTimedOut"`
	RunsErrored    int           `json:"runsErrored"`
	ElapsedMax     time.Duration `json:"elapsedMax"`
	Elapsed50th    time.Duration `json:"elapsed50th"`
	Elapsed95th    time.Duration `json:"elapsed95th"`
}
