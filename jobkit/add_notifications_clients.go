package jobkit

import (
	"github.com/blend/go-sdk/aws"
	"github.com/blend/go-sdk/aws/ses"
	"github.com/blend/go-sdk/datadog"
	"github.com/blend/go-sdk/slack"
)

// AddNotificationClients adds notification clients to a given job.
func AddNotificationClients(job *Job, cfg Config) {
	// set up myriad of notification targets
	if !cfg.AWS.IsZero() {
		job.EmailClient = ses.New(aws.MustNewSession(cfg.AWS))
	}
	if !cfg.Slack.IsZero() {
		job.SlackClient = slack.New(cfg.Slack)
	}
	if !cfg.Datadog.IsZero() {
		job.StatsClient = datadog.MustNew(cfg.Datadog)
	}
}
