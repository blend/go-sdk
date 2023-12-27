/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configutil

const (
	// EnvVarConfigPath is the env var for configs.
	EnvVarConfigPath	= "CONFIG_PATH"
	// ExtensionJSON is a file extension.
	ExtensionJSON	= ".json"
	// ExtensionYAML is a file extension.
	ExtensionYAML	= ".yaml"
	// ExtensionYML is a file extension.
	ExtensionYML	= ".yml"
)

var (
	// DefaultPaths are default path locations.
	// They are tested and read in order, so the later
	// paths will override data found in the earlier ones.
	DefaultPaths = []string{
		"/var/secrets/config.yml",
		"/var/secrets/config.yaml",
		"/var/secrets/config.json",
		"./_config/config.yml",
		"./_config/config.yaml",
		"./_config/config.json",
		"./config.yml",
		"./config.yaml",
		"./config.json",
	}
)
