package cron

import "github.com/blend/go-sdk/logger"

// JobManagerOption is a job manager option.
type JobManagerOption func(*JobManager)

// OptLog sets the job manager logger.
func OptLog(log logger.Log) JobManagerOption {
	return func(jm *JobManager) { jm.Log = log }
}

// OptTracer sets the job manager tracer.
func OptTracer(tracer Tracer) JobManagerOption {
	return func(jm *JobManager) { jm.Tracer = tracer }
}
