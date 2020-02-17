package slack

import (
	"context"

	"github.com/blend/go-sdk/configutil"
)

// Config represents the required fields for the config.
type Config struct {
	APIToken  string `json:"apiToken,omitempty" yaml:"apiToken,omitempty" env:"SLACK_TOKEN"`
	Username  string `json:"username,omitempty" yaml:"username,omitempty" env:"SLACK_USERNAME"`
	Channel   string `json:"channel,omitempty" yaml:"channel,omitempty" env:"SLACK_CHANNEL"`
	IconURL   string `json:"iconURL,omitempty" yaml:"iconURL,omitempty" env:"SLACK_ICON_URL"`
	IconEmoji string `json:"iconEmoji,omitempty" yaml:"iconEmoji,omitempty" env:"SLACK_ICON_EMOJI"`
	Webhook   string `json:"webhook,omitempty" yaml:"webhook,omitempty"  env:"SLACK_WEBHOOK"`
}

// Resolve includes extra steps on configutil.Read(...).
func (c *Config) Resolve(ctx context.Context) error {
	return configutil.GetEnvVars(ctx).ReadInto(c)
}

// IsZero returns if the config is set or not.
func (c Config) IsZero() bool {
	return len(c.Channel) == 0 && len(c.Webhook) == 0
}
