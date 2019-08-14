package jobkit

import (
	"github.com/blend/go-sdk/aws"
	"github.com/blend/go-sdk/aws/ses"
	"github.com/blend/go-sdk/datadog"
	"github.com/blend/go-sdk/email"
	"github.com/blend/go-sdk/slack"
	"github.com/blend/go-sdk/stats"
)

// AddNotificationClients adds notification clients to a given job.
func AddNotificationClients(job *Job, cfg Config) {
	// set up myriad of notification targets
	var emailClient email.Sender
	if !cfg.AWS.IsZero() {
		emailClient = ses.New(aws.MustNewSession(cfg.AWS))
	}
	var slackClient slack.Sender
	if !cfg.Slack.IsZero() {
		slackClient = slack.New(cfg.Slack)
	}
	var statsClient stats.Collector
	if !cfg.Datadog.IsZero() {
		statsClient = datadog.MustNew(cfg.Datadog)
	}

	job.WithEmailClient(emailClient).
		WithStatsClient(statsClient).
		WithSlackClient(slackClient)
}
