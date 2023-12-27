/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package profanity

// Option is a function that modifies a config.
type Option func(*Profanity)

// OptVerbose sets if we should show verbose output.
func OptVerbose(verbose bool) Option {
	return func(p *Profanity) {
		p.Config.Verbose = &verbose
	}
}

// OptDebug sets if we should show debug output.
func OptDebug(debug bool) Option {
	return func(p *Profanity) {
		p.Config.Debug = &debug
	}
}

// OptExitFirst sets if we should stop after the first failure.
func OptExitFirst(exitFirst bool) Option {
	return func(p *Profanity) {
		p.Config.ExitFirst = &exitFirst
	}
}

// OptRoot sets the root directory to start the profanity check.
func OptRoot(root string) Option {
	return func(p *Profanity) {
		p.Config.Root = root
	}
}

// OptRulesFile sets the rules file to check for in each directory.
func OptRulesFile(rulesFile string) Option {
	return func(p *Profanity) {
		p.Config.RulesFile = rulesFile
	}
}

// OptIncludeFiles sets the include glob filter for files.
func OptIncludeFiles(includeGlobs ...string) Option {
	return func(p *Profanity) {
		p.Config.Files.Filter.Include = includeGlobs
	}
}

// OptExcludeFiles sets the exclude glob filter for files.
func OptExcludeFiles(excludeGlobs ...string) Option {
	return func(p *Profanity) {
		p.Config.Files.Filter.Exclude = excludeGlobs
	}
}

// OptIncludeDirs sets the include glob filter for files.
func OptIncludeDirs(includeGlobs ...string) Option {
	return func(p *Profanity) {
		p.Config.Dirs.Filter.Include = includeGlobs
	}
}

// OptExcludeDirs sets the exclude glob filter for directories.
func OptExcludeDirs(excludeGlobs ...string) Option {
	return func(p *Profanity) {
		p.Config.Dirs.Filter.Exclude = excludeGlobs
	}
}

// OptConfig sets the config in its entirety.
func OptConfig(cfg Config) Option {
	return func(p *Profanity) {
		p.Config = cfg
	}
}
