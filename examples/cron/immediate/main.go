package main

import (
	"context"
	"os"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/logger"
)

// NOTE: Ensure that
//       * `InfrequentTask` satisfies `cron.Job`.
//       * `InfrequentTask` satisfies `cron.ScheduleProvider`.
var (
	_ cron.Job              = (*InfrequentTask)(nil)
	_ cron.ScheduleProvider = (*InfrequentTask)(nil)
)

// Config contains options for the command.
type Config struct {
	ServiceName string `json:"serviceName" yaml:"serviceName" env:"SERVICE_NAME"`
	ServiceEnv  string `json:"serviceEnv" yaml:"serviceEnv" env:"SERVICE_ENV"`
}

// NewConfigFromEnv returns a new config from the environment.
func NewConfigFromEnv() *Config {
	var config Config
	env.Env().ReadInto(&config)
	return &config
}

// InfrequentTask extends the lease on vault token.
type InfrequentTask struct {
	Config *Config
	Log    logger.Log
}

// Name returns the job name.
func (it *InfrequentTask) Name() string {
	return "infrequent_task"
}

// Schedule returns a schedule for the job.
func (it *InfrequentTask) Schedule() cron.Schedule {
	return cron.Immediately().Then(cron.EverySecond())
}

// Execute represents the job body.
func (it *InfrequentTask) Execute(ctx context.Context) error {
	sum := 0
	for i := 0; i < 10000; i++ {
		sum += i
	}
	logger.MaybeDebugf(it.Log, "Computed sum: %d", sum)
	return nil
}

func main() {
	log := logger.All()
	config := NewConfigFromEnv()
	log.Infof("starting `%s` infrequent task daemon", config.ServiceName)
	jm := cron.Default()
	cron.OptLog(log)(jm)
	job := &InfrequentTask{Config: config, Log: log}
	jm.LoadJobs(job)
	if err := graceful.Shutdown(jm); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
