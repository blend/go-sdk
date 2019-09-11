package jobkit

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/email"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/sentry"
	"github.com/blend/go-sdk/slack"
	"github.com/blend/go-sdk/stats"
	"github.com/blend/go-sdk/stringutil"
)

var (
	_ cron.Job                    = (*Job)(nil)
	_ cron.OnStartReceiver        = (*Job)(nil)
	_ cron.OnCompleteReceiver     = (*Job)(nil)
	_ cron.OnFailureReceiver      = (*Job)(nil)
	_ cron.OnCancellationReceiver = (*Job)(nil)
	_ cron.OnBrokenReceiver       = (*Job)(nil)
	_ cron.OnFixedReceiver        = (*Job)(nil)
	_ cron.OnDisabledReceiver     = (*Job)(nil)
	_ cron.OnEnabledReceiver      = (*Job)(nil)
)

// NewJob returns a new job.
func NewJob(job cron.Job, options ...JobOption) (*Job, error) {
	j := Job{
		Job: job,
	}
	var err error
	for _, opt := range options {
		if err = opt(&j); err != nil {
			return nil, err
		}
	}
	return &j, nil
}

// JobOption is an option for jobs.
type JobOption func(*Job) error

// Job is the main job body.
type Job struct {
	cron.Job
	Config       JobConfig
	Email        email.Message
	Log          logger.Log
	StatsClient  stats.Collector
	SlackClient  slack.Sender
	SentryClient sentry.Sender
	EmailClient  email.Sender
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
