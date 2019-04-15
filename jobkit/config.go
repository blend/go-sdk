package jobkit

import (
	"github.com/blend/go-sdk/airbrake"
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
	MaxLogBytes int             `yaml:"maxLogBytes"`
	Cron        cron.Config     `yaml:"cron"`
	Logger      logger.Config   `yaml:"logger"`
	Web         web.Config      `yaml:"web"`
	Airbrake    airbrake.Config `yaml:"airbrake"`
	AWS         aws.Config      `yaml:"aws"`
	Email       email.Message   `yaml:"email"`
	Datadog     datadog.Config  `yaml:"datadog"`
	Slack       slack.Config    `yaml:"slack"`
}

// Resolve applies resolution steps to the config.
func (c *Config) Resolve() error {
	return configutil.AnyError(
		c.Cron.Resolve(),
		c.Logger.Resolve(),
		c.Web.Resolve(),
		c.Airbrake.Resolve(),
		c.AWS.Resolve(),
		c.Email.Resolve(),
		c.Datadog.Resolve(),
		c.Slack.Resolve(),
	)
}

// MaxLogBytesOrDefault is a the maximum amount of log data to buffer.
func (c Config) MaxLogBytesOrDefault() int {
	if c.MaxLogBytes > 0 {
		return c.MaxLogBytes
	}
	return DefaultMaxLogBytes
}
