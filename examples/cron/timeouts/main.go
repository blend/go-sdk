package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/logger"
)

type emptyJob struct {
	running bool
}

func (j *emptyJob) Timeout() time.Duration {
	return 2 * time.Second
}

func (j *emptyJob) Name() string {
	return "printJob"
}

func (j *emptyJob) Execute(ctx context.Context) error {
	j.running = true
	var runFor = 8 * time.Second
	if rand.Int()%2 == 1 {
		runFor = time.Second
	}

	alarm := time.After(runFor)
	select {
	case <-alarm:
		j.running = false
		return nil
	case <-ctx.Done():
		j.running = false
		return nil
	}
}

func (j *emptyJob) OnCancellation(_ *cron.JobInvocation) {
	j.running = false
}

func (j *emptyJob) Status() string {
	if j.running {
		return "Request in progress"
	}
	return "Request idle."
}

func (j *emptyJob) Schedule() cron.Schedule {
	return cron.Immediately().Then(cron.Every(10 * time.Second))
}

func main() {
	jm := cron.New(cron.OptLog(logger.All()))
	jm.LoadJobs(&emptyJob{})
	if err := jm.StartAsync(); err != nil {
		logger.FatalExit(err)
	}

	for {
		status := jm.Status()
		for _, job := range status.Jobs {
			if runningJob, ok := status.Running[job.Name]; ok {
				jm.Log.Infof("job: %s > %s state: running elapsed: %v", job.Name, runningJob.ID, cron.Since(runningJob.Started))
			} else {
				jm.Log.Infof("job: %s state: stopped", job.Name)
			}
		}

		time.Sleep(1000 * time.Millisecond)
	}
}
