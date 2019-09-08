package cron

import (
	"time"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/ref"
)

// Config governs the base options all jobs inherit from.
type Config struct {
	HistoryEnabled  *bool         `json:"historyEnabled" yaml:"historyEnabled" env:"CRON_HISTORY_ENABLED"`
	HistoryMaxCount int           `json:"historyMaxCount" yaml:"historyMaxCount" env:"CRON_HISTORY_MAX_COUNT"`
	HistoryMaxAge   time.Duration `json:"historyMaxAge" yaml:"historyMaxAge" env:"CRON_HISTORY_MAX_AGE"`
}

// Resolve adds extra resolution steps when reading the config.
func (hc Config) Resolve() error {
	return configutil.AnyError(
		configutil.SetBool(&hc.HistoryEnabled, configutil.Bool(hc.HistoryEnabled), configutil.Parse(configutil.Env("CRON_HISTORY_ENABLED")), configutil.Bool(ref.Bool(DefaultHistoryEnabled))),
		configutil.SetInt(&hc.HistoryMaxCount, configutil.Int(hc.HistoryMaxCount), configutil.Parse(configutil.Env("CRON_HISTORY_MAX_COUNT")), configutil.Int(DefaultHistoryMaxCount)),
		configutil.SetDuration(&hc.HistoryMaxAge, configutil.Duration(hc.HistoryMaxAge), configutil.Parse(configutil.Env("CRON_HISTORY_MAX_AGE")), configutil.Duration(DefaultHistoryMaxAge)),
	)
}

// HistoryEnabledOrDefault returns if history is enabled.
func (hc Config) HistoryEnabledOrDefault() bool {
	if hc.HistoryEnabled != nil {
		return *hc.HistoryEnabled
	}
	return DefaultHistoryEnabled
}

// HistoryMaxCountOrDefault returns the max count or a default.
func (hc Config) HistoryMaxCountOrDefault() int {
	if hc.HistoryMaxCount > 0 {
		return hc.HistoryMaxCount
	}
	return DefaultHistoryMaxCount
}

// HistoryMaxAgeOrDefault returns the max age or a default.
func (hc Config) HistoryMaxAgeOrDefault() time.Duration {
	if hc.HistoryMaxAge > 0 {
		return hc.HistoryMaxAge
	}
	return DefaultHistoryMaxAge
}
