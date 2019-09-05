package main

import (
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/blend/go-sdk/aws"
	"github.com/blend/go-sdk/aws/ses"
	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/datadog"
	"github.com/blend/go-sdk/email"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/graceful"
	"github.com/blend/go-sdk/jobkit"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/ref"
	"github.com/blend/go-sdk/slack"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/stringutil"
	"github.com/blend/go-sdk/webutil"
)

var (
	flagBind                    *string
	flagConfigPath              *string
	flagDefaultJobName          *string
	flagDefaultJobExec          *string
	flagDefaultJobSchedule      *string
	flagDefaultJobTimeout       *time.Duration
	flagDefaultJobDiscardOutput *bool
	flagDefaultJobSerial        *bool
	flagDisableServer           *bool
)

type config struct {
	jobkit.Config `yaml:",inline"`
	DisableServer *bool       `yaml:"disableServer"`
	Jobs          []jobConfig `yaml:"jobs"`
}

func (c *config) Resolve() error {
	if err := configutil.SetString(&c.Web.BindAddr, configutil.String(*flagBind), configutil.Env("BIND_ADDR"), configutil.String(c.Web.BindAddr)); err != nil {
		return err
	}
	if err := configutil.SetBool(&c.DisableServer, configutil.Bool(flagDisableServer), configutil.Bool(c.DisableServer), configutil.Bool(ref.Bool(false))); err != nil {
		return err
	}
	return nil
}

type jobConfig struct {
	// JobConfig is the default jobkit job options.
	jobkit.JobConfig `yaml:",inline"`
	// Exec is the command to execute.
	Exec []string `yaml:"exec"`
	// DiscardOutput indicates if we should discard output.
	DiscardOutput *bool `yaml:"discardOutput"`
}

// DiscardOutputOrDefault returns a value or a default.
func (jc *jobConfig) DiscardOutputOrDefault() bool {
	if jc.DiscardOutput != nil {
		return *jc.DiscardOutput
	}
	return false
}

func (jc *jobConfig) Resolve() error {
	return configutil.AnyError(
		configutil.SetString(&jc.Name, configutil.String(*flagDefaultJobName), configutil.String(env.Env().ServiceName()), configutil.String(jc.Name), configutil.String(stringutil.Letters.Random(8))),
		configutil.SetBool(&jc.DiscardOutput, configutil.Bool(flagDefaultJobDiscardOutput), configutil.Bool(jc.DiscardOutput), configutil.Bool(ref.Bool(false))),
		configutil.SetBool(&jc.Serial, configutil.Bool(flagDefaultJobSerial), configutil.Bool(jc.Serial), configutil.Bool(ref.Bool(true))),
		configutil.SetString(&jc.Schedule, configutil.String(*flagDefaultJobSchedule), configutil.String(jc.Schedule)),
		configutil.SetDuration(&jc.Timeout, configutil.Duration(*flagDefaultJobTimeout), configutil.Duration(jc.Timeout)),
	)
}

func command() *cobra.Command {
	return &cobra.Command{
		Use:   "job",
		Short: "Job runs a command on a schedule, and tracks limited job history in memory.",
		Long:  "Job runs a command on a schedule, and tracks limited job history in memory.",
		Example: `
# echo 'hello world' with the default schedule
job -- echo 'hello world'

# echo 'hello world' every 30 seconds
job --schedule='*/30 * * * *' -- echo 'hello world'

# set the job name
job -n echo --schedule='*/30 * * * *' -- echo 'hello world'

# use a config
job -c config.yml'

# where the config can specify multiple jobs.
# it can also set general configuration options like the bind address (located in the web config).
"""
logger:
  flags: [all, -http.request]

web:
  bindAddr: :8080

jobs:
- name: echo
  schedule: '*/30 * * * *'
  exec: [echo, 'hello world']
- name: echo2
  schedule: '*/30 * * * *'
  exec: [echo, 'hello again']
"""
`,
	}
}

func main() {
	cmd := command()
	cmd.Run = fatalExit(run)

	flagBind = cmd.Flags().String("bind", "", "The management http server bind address.")
	flagConfigPath = cmd.Flags().StringP("config", "c", "", "The config path.")
	flagDefaultJobName = cmd.Flags().StringP("name", "n", "", "The job name (will default to a random string of 8 letters).")
	flagDefaultJobSchedule = cmd.Flags().StringP("schedule", "s", "", "The job schedule in cron format (ex: '*/5 * * * *')")
	flagDefaultJobTimeout = cmd.Flags().Duration("timeout", 0, "The job execution timeout as a duration (ex: 5s)")
	flagDefaultJobDiscardOutput = cmd.Flags().Bool("discard-output", false, "If jobs should discard console output from the action.")
	flagDefaultJobSerial = cmd.Flags().Bool("serial", true, "The job should run serially (that is, only one can be active at a time)")
	flagDisableServer = cmd.Flags().Bool("disable-server", false, "If the management server should be disabled.")

	if err := cmd.Execute(); err != nil {
		logger.FatalExit(err)
	}
}

func fatalExit(action func(*cobra.Command, []string) error) func(*cobra.Command, []string) {
	return func(parent *cobra.Command, args []string) {
		if err := action(parent, args); err != nil {
			logger.FatalExit(err)
		}
	}
}

func run(cmd *cobra.Command, args []string) error {
	var cfg config
	if _, err := configutil.Read(&cfg, configutil.OptPaths(*flagConfigPath)); !configutil.IsIgnored(err) {
		return err
	}

	log, err := logger.New(logger.OptConfig(cfg.Logger))
	if err != nil {
		return err
	}
	log.Flags.Enable(
		logger.Error,
		logger.Fatal,
		cron.FlagStarted,
		cron.FlagComplete,
		cron.FlagFixed,
		cron.FlagBroken,
		cron.FlagFailed,
		cron.FlagCancelled,
		webutil.HTTPRequest,
		webutil.HTTPResponse,
	)

	log.Debugf("using logger flags: %s", log.Flags.String())

	defaultJobCfg, err := createDefaultJobConfig(args...)
	if err != nil {
		return err
	}
	if defaultJobCfg != nil {
		cfg.Jobs = append(cfg.Jobs, *defaultJobCfg)
	}

	if len(cfg.Jobs) == 0 {
		return ex.New("must supply a command to run with `--exec=...` or `-- command`), or provide a jobs config file")
	}

	// set up myriad of notification targets
	var emailClient email.Sender
	if !cfg.AWS.IsZero() {
		emailClient = ses.New(aws.MustNewSession(cfg.AWS))
		log.Infof("adding ses email notifications")
	} else if !cfg.SMTP.IsZero() {
		emailClient = cfg.SMTP
		log.Infof("adding smtp email notifications")
	}

	if !cfg.EmailDefaults.IsZero() {
		log.Debugf("using email defaults; from=%s", cfg.EmailDefaults.From)
		log.Debugf("using email defaults; to=%s", stringutil.CSV(cfg.EmailDefaults.To))
	}

	var slackClient slack.Sender
	if !cfg.Slack.IsZero() {
		slackClient = slack.New(cfg.Slack)
		log.Infof("adding slack notifications")
	}
	var statsClient stats.Collector
	if !cfg.Datadog.IsZero() {
		statsClient, err = datadog.New(cfg.Datadog)
		if err != nil {
			return err
		}
		log.Infof("adding datadog metrics")
	}

	jobs := cron.New(cron.OptConfig(cfg.Config.Cron), cron.OptLog(log.WithPath("cron")))

	for _, jobCfg := range cfg.Jobs {
		job, err := createJobFromConfig(cfg, jobCfg)
		if err != nil {
			return err
		}

		job.Log = log
		job.EmailClient = emailClient
		job.SlackClient = slackClient
		job.StatsClient = statsClient

		var serial string
		if jobCfg.SerialOrDefault() {
			serial = "serial execution"
		} else {
			serial = "parallel execution"
		}

		log.Infof("loading job `%s` with schedule `%s` with %v", jobCfg.Name, jobCfg.ScheduleOrDefault(), serial)
		jobs.LoadJobs(job)
	}

	hosted := []graceful.Graceful{jobs}

	if !*flagDisableServer {
		ws := jobkit.NewManagementServer(jobs, cfg.Config)
		ws.Log = log.WithPath("management server")
		hosted = append(hosted, ws)
	} else {
		log.Infof("management server disabled")
	}
	return graceful.Shutdown(hosted...)
}

func createDefaultJobConfig(args ...string) (*jobConfig, error) {
	cfg := new(jobConfig)
	if err := cfg.Resolve(); err != nil {
		return nil, err
	}
	cfg.Exec = args
	if len(cfg.Exec) == 0 {
		return nil, nil
	}
	return cfg, nil
}

func createJobFromConfig(base config, cfg jobConfig) (*jobkit.Job, error) {
	if len(cfg.Exec) == 0 {
		return nil, ex.New("job exec and command unset", ex.OptMessagef("job: %s", cfg.Name))
	}
	action := jobkit.CreateShellAction(cfg.Exec, jobkit.OptShellActionDiscardOutput(cfg.DiscardOutputOrDefault()))
	job, err := jobkit.NewJob(cfg.JobConfig, action)
	if err != nil {
		return nil, err
	}
	if job.Config.Description == "" {
		job.Config.Description = strings.Join(cfg.Exec, " ")
	}
	job.EmailDefaults = email.MergeMessages(base.EmailDefaults, cfg.EmailDefaults)
	return job, nil
}
