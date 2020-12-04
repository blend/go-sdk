package profanity

// Config is the profanity rules parsing config.
type Config struct {
	Verbose   *bool      `yaml:"verbose"`
	Debug     *bool      `yaml:"debug"`
	FailFast  *bool      `yaml:"failFast"`
	Path      string     `yaml:"path"`
	RulesFile string     `yaml:"rulesFile"`
	Files     GlobFilter `yaml:"files"`
	Dirs      GlobFilter `yaml:"dirs"`
}

// PathOrDefault returns the starting path or a default.
func (c Config) PathOrDefault() string {
	if c.Path != "" {
		return c.Path
	}
	return DefaultPath
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
