package profanity

// Config is the profanity rules parsing config.
type Config struct {
	Verbose   *bool    `yaml:"verbose"`
	Debug     *bool    `yaml:"debug"`
	FailFast  *bool    `yaml:"failFast"`
	Root      string   `yaml:"root"`
	RulesFile string   `yaml:"rulesFile"`
	Include   []string `yaml:"include,omitempty"`
	Exclude   []string `yaml:"exclude,omitempty"`
}

// VerboseOrDefault returns an option or a default.
func (c Config) VerboseOrDefault() bool {
	if c.Verbose != nil {
		return *c.Verbose
	}
	return false
}

// DebugOrDefault returns an option or a default.
func (c Config) DebugOrDefault() bool {
	if c.Debug != nil {
		return *c.Debug
	}
	return false
}

// FailFastOrDefault returns an option or a default.
func (c Config) FailFastOrDefault() bool {
	if c.FailFast != nil {
		return *c.FailFast
	}
	return false
}

// RulesFileOrDefault returns the rules file or a default.
func (c Config) RulesFileOrDefault() string {
	if c.RulesFile != "" {
		return c.RulesFile
	}
	return DefaultRulesFile
}
