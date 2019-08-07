package profanity

import "github.com/blend/go-sdk/ref"

// ConfigOption is a function that modifies a config.
type ConfigOption func(*Config)

// OptVerbose sets if we should show verbose output.
func OptVerbose(verbose bool) ConfigOption {
	return func(c *Config) {
		c.Verbose = ref.Bool(verbose)
	}
}

// OptDebug sets if we should show debug output.
func OptDebug(debug bool) ConfigOption {
	return func(c *Config) {
		c.Debug = ref.Bool(debug)
	}
}

// OptFailFast sets if we should stop after the first failure.
func OptFailFast(failFast bool) ConfigOption {
	return func(c *Config) {
		c.FailFast = ref.Bool(failFast)
	}
}

// OptRoot sets the root directory to start the profanity check.
func OptRoot(root string) ConfigOption {
	return func(c *Config) {
		c.Root = root
	}
}

// OptRulesFile sets the rules file to check for in each directory.
func OptRulesFile(rulesFile string) ConfigOption {
	return func(c *Config) {
		c.RulesFile = rulesFile
	}
}

// OptInclude sets the include filter.
func OptInclude(includes ...string) ConfigOption {
	return func(c *Config) {
		c.Include = includes
	}
}

// OptExclude sets the exclude filter.
func OptExclude(excludes ...string) ConfigOption {
	return func(c *Config) {
		c.Exclude = excludes
	}
}

// OptConfig sets the config in its entirety.
func OptConfig(cfg Config) ConfigOption {
	return func(c *Config) {
		*c = cfg
	}
}
