/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package codeowners

// Config are the configuration options for the utility.
type Config struct {
	// Path is the owners file path (typically `.github/CODEOWNERS`)
	Path	string	`yaml:"path"`
	// GithubURL is the url of the github instance to communicate with.
	GithubURL	string	`yaml:"githubURL"`
	// GithubToken is the authorization token used to communicate with github.
	GithubToken	string	`yaml:"githubToken"`

	// Quiet controls whether output is suppressed.
	Quiet	*bool	`yaml:"quiet"`
	// Verbose controls whether verbose output is shown.
	Verbose	*bool	`yaml:"verbose"`
	// Debug controls whether debug output is shown.
	Debug	*bool	`yaml:"debug"`
}

// PathOrDefault is the path for the codeowners file or a default.
func (c Config) PathOrDefault() string {
	if c.Path != "" {
		return c.Path
	}
	return DefaultPath
}

// GithubURLOrDefault returns a value or a default.
func (c Config) GithubURLOrDefault() string {
	if c.GithubURL != "" {
		return c.GithubURL
	}
	return DefaultGithubURL
}

// QuietOrDefault returns a value or a default.
func (c Config) QuietOrDefault() bool {
	if c.Quiet != nil {
		return *c.Quiet
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
