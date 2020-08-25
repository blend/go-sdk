package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/ref"
)

var (
	flagTarget      = flag.String("target", "", "The target URL")
	flagEnvironment = flag.String("env", "", "The current environment")
)

// Config is a sample config.
type Config struct {
	Target       string `yaml:"target"`
	DebugEnabled *bool  `yaml:"debugEnabled"`
	MaxCount     int    `yaml:"maxCount"`
	Environment  string `yaml:"env"`
}

// Resolve parses the config and sets values from predefined sources.
func (c *Config) Resolve(ctx context.Context) error {
	return configutil.Resolve(ctx,
		configutil.SetString(&c.Target,
			configutil.String(*flagTarget),
			configutil.Env("TARGET"),
			configutil.String(c.Target),
			configutil.String("https://google.com/robots.txt"),
		),
		configutil.SetBool(&c.DebugEnabled,
			configutil.Env("DEBUG_ENABLED"),
			configutil.Bool(c.DebugEnabled),
			configutil.Bool(ref.Bool(true)),
		),
		configutil.SetInt(&c.MaxCount,
			configutil.Int(c.MaxCount),
			configutil.Env("MAX_COUNT"),
			configutil.Int(10),
		),
		configutil.SetString(&c.Environment,
			configutil.String(*flagEnvironment),
			configutil.Env("SERVICE_ENV"),
			configutil.String(c.Environment),
			configutil.String("development"),
		),
	)
}

var (
	_ configutil.Resolver = (*Config)(nil)
)

func main() {
	flag.Parse()
	config := new(Config)
	if _, err := configutil.Read(config, configutil.OptLog(logger.All().WithPath("config"))); !configutil.IsIgnored(err) {
		logger.FatalExit(err)
	}
	fmt.Println("target:", config.Target)
	fmt.Println("debug enabled:", *config.DebugEnabled)
	fmt.Println("env:", config.Environment)
}
