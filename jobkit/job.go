package jobkit

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/email"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/ref"
	"github.com/blend/go-sdk/sentry"
	"github.com/blend/go-sdk/slack"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/stringutil"
)

var (
	_ cron.Job                    = (*Job)(nil)
	_ cron.LabelsProvider         = (*Job)(nil)
	_ cron.JobConfigProvider      = (*Job)(nil)
	_ cron.ScheduleProvider       = (*Job)(nil)
	_ cron.OnStartReceiver        = (*Job)(nil)
	_ cron.OnCompleteReceiver     = (*Job)(nil)
	_ cron.OnFailureReceiver      = (*Job)(nil)
	_ cron.OnCancellationReceiver = (*Job)(nil)
	_ cron.OnBrokenReceiver       = (*Job)(nil)
	_ cron.OnFixedReceiver        = (*Job)(nil)
	_ cron.OnDisabledReceiver     = (*Job)(nil)
	_ cron.OnEnabledReceiver      = (*Job)(nil)
	_ cron.HistoryProvider        = (*Job)(nil)
)

// NewJob returns a new job.
func NewJob(cfg JobConfig, action func(context.Context) error, options ...JobOption) (*Job, error) {
	options = append([]JobOption{
		OptConfig(cfg),
		OptAction(action),
		OptParsedSchedule(cfg.ScheduleOrDefault()),
	}, options...)

	var job Job
	var err error
	for _, opt := range options {
		if err = opt(&job); err != nil {
			return nil, err
		}
	}
	return &job, nil
}

// Wrap wraps a cron job with the jobkit notifications.
func Wrap(job cron.Job) *Job {
	j := &Job{
		Config: JobConfig{
			JobConfig: cron.JobConfig{
				Name: job.Name(),
			},
		},
		Action: job.Execute,
	}

	if typed, ok := job.(cron.ScheduleProvider); ok {
		j.CompiledSchedule = typed.Schedule()
	}
	if typed, ok := job.(cron.DescriptionProvider); ok {
		j.Config.Description = typed.Description()
	}
	if typed, ok := job.(cron.LabelsProvider); ok {
		j.Config.Labels = typed.Labels()
	}
	if typed, ok := job.(cron.TimeoutProvider); ok {
		j.Config.Timeout = typed.Timeout()
	}
	if typed, ok := job.(cron.ShutdownGracePeriodProvider); ok {
		j.Config.ShutdownGracePeriod = typed.ShutdownGracePeriod()
	}
	if typed, ok := job.(cron.HistoryDisabledProvider); ok {
		j.Config.HistoryDisabled = ref.Bool(typed.HistoryDisabled())
	}
	if typed, ok := job.(cron.HistoryMaxCountProvider); ok {
		j.Config.HistoryMaxCount = typed.HistoryMaxCount()
	}
	if typed, ok := job.(cron.HistoryMaxAgeProvider); ok {
		j.Config.HistoryMaxAge = typed.HistoryMaxAge()
	}
	if typed, ok := job.(cron.ShouldSkipLoggerListenersProvider); ok {
		j.Config.ShouldSkipLoggerListeners = ref.Bool(typed.ShouldSkipLoggerListeners())
	}
	if typed, ok := job.(cron.ShouldSkipLoggerOutputProvider); ok {
		j.Config.ShouldSkipLoggerOutput = ref.Bool(typed.ShouldSkipLoggerOutput())
	}
	return j
}

// OptAction sets the job action.
func OptAction(action func(context.Context) error) JobOption {
	return func(job *Job) error {
		job.Action = action
		return nil
	}
}

// OptConfig sets the job config.
func OptConfig(cfg JobConfig) JobOption {
	return func(job *Job) error {
		job.Config = cfg
		return nil
	}
}

// OptParsedSchedule sets the job's compiled schedule from a schedule string.
func OptParsedSchedule(schedule string) JobOption {
	return func(job *Job) error {
		schedule, err := cron.ParseString(schedule)
		if err != nil {
			return err
		}
		job.CompiledSchedule = schedule
		return nil
	}
}

// JobOption is an option for jobs.
type JobOption func(*Job) error

// Job is the main job body.
type Job struct {
	Config JobConfig

	CompiledSchedule cron.Schedule
	Email            email.Message
	Action           func(context.Context) error

	Log          logger.Log
	StatsClient  stats.Collector
	SlackClient  slack.Sender
	SentryClient sentry.Sender
	EmailClient  email.Sender
}

// Name returns the job name.
func (job Job) Name() string {
	return job.Config.Name
}

// Labels returns the job labels.
func (job Job) Labels() map[string]string {
	return job.Config.Labels
}

// Schedule returns the job schedule.
func (job Job) Schedule() cron.Schedule {
	return job.CompiledSchedule
}

// JobConfig returns the underlying job config.
func (job Job) JobConfig() cron.JobConfig {
	return job.Config.JobConfig
}

// OnStart is a lifecycle event handler.
func (job Job) OnStart(ctx context.Context) {
	job.stats(ctx, cron.FlagStarted)
	if job.Config.NotifyOnStartOrDefault() {
		job.notify(ctx, cron.FlagStarted)
	}
}

// OnComplete is a lifecycle event handler.
func (job Job) OnComplete(ctx context.Context) {
	job.stats(ctx, cron.FlagComplete)
	if job.Config.NotifyOnSuccessOrDefault() {
		job.notify(ctx, cron.FlagComplete)
	}
}

// OnFailure is a lifecycle event handler.
func (job Job) OnFailure(ctx context.Context) {
	job.stats(ctx, cron.FlagFailed)
	if job.Config.NotifyOnFailureOrDefault() {
		job.notify(ctx, cron.FlagFailed)
	}
}

// OnBroken is a lifecycle event handler.
func (job Job) OnBroken(ctx context.Context) {
	job.stats(ctx, cron.FlagBroken)
	if job.Config.NotifyOnBrokenOrDefault() {
		job.notify(ctx, cron.FlagBroken)
	}
}

// OnFixed is a lifecycle event handler.
func (job Job) OnFixed(ctx context.Context) {
	job.stats(ctx, cron.FlagFixed)
	if job.Config.NotifyOnFixedOrDefault() {
		job.notify(ctx, cron.FlagFixed)
	}
}

// OnCancellation is a lifecycle event handler.
func (job Job) OnCancellation(ctx context.Context) {
	job.stats(ctx, cron.FlagCancelled)
	if job.Config.NotifyOnCancellationOrDefault() {
		job.notify(ctx, cron.FlagCancelled)
	}
}

// OnEnabled is a lifecycle event handler.
func (job Job) OnEnabled(ctx context.Context) {
	if job.Config.NotifyOnEnabledOrDefault() {
		job.notify(ctx, cron.FlagEnabled)
	}
}

// OnDisabled is a lifecycle event handler.
func (job Job) OnDisabled(ctx context.Context) {
	if job.Config.NotifyOnDisabledOrDefault() {
		job.notify(ctx, cron.FlagDisabled)
	}
}

// PersistHistory writes the history to disk.
// It does so completely.
func (job Job) PersistHistory(ctx context.Context, log []cron.JobInvocation) error {
	historyDirectory := job.Config.HistoryPathOrDefault()
	if _, err := os.Stat(historyDirectory); err != nil {
		if err := os.MkdirAll(historyDirectory, 0755); err != nil {
			return ex.New(err)
		}
	}
	historyPath := filepath.Join(historyDirectory, stringutil.Slugify(job.Name())+".json")
	f, err := os.Create(historyPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(log)
}

// RestoreHistory restores history from disc.
func (job Job) RestoreHistory(ctx context.Context) (output []cron.JobInvocation, err error) {
	historyPath := filepath.Join(job.Config.HistoryPathOrDefault(), stringutil.Slugify(job.Name())+".json")
	if _, statErr := os.Stat(historyPath); statErr != nil {
		return
	}
	var f *os.File
	f, err = os.Open(historyPath)
	if err != nil {
		return
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&output)
	return
}

// Execute is the job body.
func (job Job) Execute(ctx context.Context) error {
	return job.Action(ctx)
}

//
// exported utility methods
//

// Debugf logs a debug message if the logger is set.
func (job Job) Debugf(ctx context.Context, format string, args ...interface{}) {
	if job.Log != nil {
		job.Log.WithPath(job.Name(), cron.GetJobInvocation(ctx).ID).WithContext(ctx).Debugf(format, args...)
	}
}

// Error logs an error if the logger i set.
func (job Job) Error(ctx context.Context, err error) error {
	if job.Log != nil {
		job.Log.WithPath(job.Name(), cron.GetJobInvocation(ctx).ID).WithContext(ctx).Error(err)
	}
	return err

}

//
// private utility methods
//

func (job Job) stats(ctx context.Context, flag string) {
	if job.StatsClient != nil {
		job.StatsClient.Increment(string(flag), fmt.Sprintf("%s:%s", stats.TagJob, job.Name()))
		if ji := cron.GetJobInvocation(ctx); ji != nil {
			job.Error(ctx, job.StatsClient.TimeInMilliseconds(string(flag), ji.Elapsed, fmt.Sprintf("%s:%s", stats.TagJob, job.Name())))
		}
	} else {
		job.Debugf(ctx, "stats client unset, skipping logging stats")
	}
}

func (job Job) notify(ctx context.Context, flag string) {
	if job.SlackClient != nil {
		if ji := cron.GetJobInvocation(ctx); ji != nil {
			job.Error(ctx, job.SlackClient.Send(context.Background(), NewSlackMessage(ji)))
		}
	} else {
		job.Debugf(ctx, "notify (slack); sender unset skipping sending slack notification")
	}

	if job.EmailClient != nil {
		if ji := cron.GetJobInvocation(ctx); ji != nil {
			message, err := NewEmailMessage(job.Email, ji)
			if err != nil {
				job.Error(ctx, err)
			}
			job.Error(ctx, job.EmailClient.Send(context.Background(), message))
			job.Debugf(ctx, "notify (email); sent email notification to %s (%s)", stringutil.CSV(message.To), message.Subject)
		} else {
			job.Debugf(ctx, "notify (email); job invocation not found on context")
		}
	} else {
		job.Debugf(ctx, "notify (email); email sender unset, skipping sending email notification")
	}
}
