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

// OptPath sets the root directory to start the profanity check.
func OptPath(path string) ConfigOption {
	return func(c *Config) {
		c.Path = path
	}
}

// OptRulesFile sets the rules file to check for in each directory.
func OptRulesFile(rulesFile string) ConfigOption {
	return func(c *Config) {
		c.RulesFile = rulesFile
	}
}

// OptFilesInclude sets the include glob filter for files.
func OptFilesInclude(includeGlobs ...string) ConfigOption {
	return func(c *Config) {
		c.Files.Filter.Include = includeGlobs
	}
}

// OptFilesExclude sets the exclude glob filter for files.
func OptFilesExclude(excludeGlobs ...string) ConfigOption {
	return func(c *Config) {
		c.Files.Filter.Exclude = excludeGlobs
	}
}

// OptDirsInclude sets the include glob filter for files.
func OptDirsInclude(includeGlobs ...string) ConfigOption {
	return func(c *Config) {
		c.Dirs.Filter.Include = includeGlobs
	}
}

// OptDirsExclude sets the exclude glob filter for directories.
func OptDirsExclude(excludeGlobs ...string) ConfigOption {
	return func(c *Config) {
		c.Dirs.Filter.Exclude = excludeGlobs
	}
}

// OptConfig sets the config in its entirety.
func OptConfig(cfg Config) ConfigOption {
	return func(c *Config) {
		*c = cfg
	}
}
