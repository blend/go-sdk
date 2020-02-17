package main

import (
	"context"
	"flag"
	"fmt"

	cfg "github.com/blend/go-sdk/configutil"
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
	return cfg.ReturnFirst(
		cfg.SetString(&c.Target, cfg.String(*flagTarget), cfg.Env(ctx, "TARGET"), cfg.String(c.Target), cfg.String("https://google.com/robots.txt")),
		cfg.SetBool(&c.DebugEnabled, cfg.Env(ctx, "DEBUG_ENABLED"), cfg.Bool(c.DebugEnabled), cfg.Bool(ref.Bool(true))),
		cfg.SetInt(&c.MaxCount, cfg.Int(c.MaxCount), cfg.Env(ctx, "MAX_COUNT"), cfg.Int(10)),
		cfg.SetString(&c.Environment, cfg.String(*flagEnvironment), cfg.Env(ctx, "SERVICE_ENV"), cfg.String(c.Environment), cfg.String("development")),
	)
}

var (
	_ cfg.Resolver = (*Config)(nil)
)

func main() {
	flag.Parse()
	config := new(Config)
	if _, err := cfg.Read(config, cfg.OptLog(logger.All().WithPath("config"))); !cfg.IsIgnored(err) {
		logger.FatalExit(err)
	}
	fmt.Println("target:", config.Target)
	fmt.Println("debug enabled:", *config.DebugEnabled)
	fmt.Println("env:", config.Environment)
}
