/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package profanity

// Config is the profanity rules parsing config.
type Config struct {
	Root		string		`yaml:"root"`
	ExitFirst	*bool		`yaml:"failFast"`
	RulesFile	string		`yaml:"rulesFile"`
	Rules		GlobFilter	`yaml:"rules"`
	Files		GlobFilter	`yaml:"files"`
	Dirs		GlobFilter	`yaml:"dirs"`
	Verbose		*bool		`yaml:"verbose"`
	Debug		*bool		`yaml:"debug"`
}

// RootOrDefault returns the starting path or a default.
func (c Config) RootOrDefault() string {
	if c.Root != "" {
		return c.Root
	}
	return DefaultRoot
}

// RulesFileOrDefault returns the rules file or a default.
func (c Config) RulesFileOrDefault() string {
	if c.RulesFile != "" {
		return c.RulesFile
	}
	return DefaultRulesFile
}

// ExitFirstOrDefault returns a value or a default.
func (c Config) ExitFirstOrDefault() bool {
	if c.ExitFirst != nil {
		return *c.ExitFirst
	}
	return false
}

// VerboseOrDefault returns a value or a default.
func (c Config) VerboseOrDefault() bool {
	if c.Verbose != nil {
		return *c.Verbose
	}
	return false
}

// DebugOrDefault returns a value or a default.
func (c Config) DebugOrDefault() bool {
	if c.Debug != nil {
		return *c.Debug
	}
	return false
}
