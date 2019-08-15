package airbrake

import (
	"github.com/blend/go-sdk/configutil"
)

// Config is the airbrake config.
type Config struct {
	ProjectID   string `json:"projectID" yaml:"projectID" env:"AIRBRAKE_PROJECT_ID"`
	ProjectKey  string `json:"projectKey" yaml:"projectKey" env:"AIRBRAKE_PROJECT_KEY"`
	Environment string `json:"environment" yaml:"environment" env:"SERVICE_ENV"`
}

// Resolve resolves config defaults.
func (c *Config) Resolve() error {
	return configutil.AnyError(
		configutil.SetString(&c.ProjectID, configutil.String(c.ProjectID), configutil.Env("AIRBRAKE_PROJECT_ID")),
		configutil.SetString(&c.ProjectKey, configutil.String(c.ProjectKey), configutil.Env("AIRBRAKE_PROJECT_KEY")),
		configutil.SetString(&c.Environment, configutil.String(c.Environment), configutil.Env("SERVICE_ENV")),
	)
}

// IsZero returns if the config is set or not.
func (c Config) IsZero() bool {
	return len(c.ProjectKey) == 0 || len(c.ProjectID) == 0
}
