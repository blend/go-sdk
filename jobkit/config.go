package jobkit

import (
	"github.com/blend/go-sdk/aws"
	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/datadog"
	"github.com/blend/go-sdk/email"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/slack"
	"github.com/blend/go-sdk/web"
)

// Config is the jobkit config.
type Config struct {
	// HistoryPath is the default historyPath.
	HistoryPath string `yaml:"historyPath"`
	// Cron is the cron manager config.
	Cron cron.Config `yaml:"cron"`
	// Logger is the logger config.
	Logger logger.Config `yaml:"logger"`
	// Web is the web config used for the management server.
	Web web.Config `yaml:"web"`
	// EmailDefaults sets the base options for notification emails.
	EmailDefaults email.Message `yaml:"emailDefaults"`
	// AWS is used by aws options like SES.
	AWS aws.Config `yaml:"aws"`
	// SMTP is the smtp options.
	SMTP email.SMTPSender `yaml:"smtp"`
	// Datadog configures the datadog client.
	Datadog datadog.Config `yaml:"datadog"`
	// Slack configues the slack webhook sender.
	Slack slack.Config `yaml:"slack"`
}

// Resolve applies resolution steps to the config.
func (c *Config) Resolve() error {
	return configutil.AnyError(
		c.Cron.Resolve(),
		c.Logger.Resolve(),
		c.Web.Resolve(),
		c.AWS.Resolve(),
		c.EmailDefaults.Resolve(),
		c.Datadog.Resolve(),
		c.Slack.Resolve(),
	)
}
