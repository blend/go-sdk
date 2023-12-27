/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package configmeta

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
)

// EnvVars
const (
	EnvVarProjectConfigPath	= "PROJECT_CONFIG_PATH"
	EnvVarConfigPath	= "CONFIG_PATH"
)

// Defaults
const (
	DefaultConfigFileName = "config.yml"
)

// Paths returns a configutil option that adds known default config locations.
func Paths() configutil.Option {
	return PathsFile(DefaultConfigFileName)
}

// PathsFile returns a configutil option that adds known default config locations.
func PathsFile(filename string) configutil.Option {
	return PathsFileContext(context.Background(), filename)
}

// PathsFileContext returns a configutil option that adds known default config locations.
func PathsFileContext(ctx context.Context, filename string) configutil.Option {
	projectName := env.GetVars(ctx).String(env.VarProjectName, env.GetVars(ctx).ServiceName())
	fallbackProjectConfigPath := fmt.Sprintf("/var/secrets/projects/%s/%s", projectName, filename)
	fallbackConfigPath := fmt.Sprintf("/var/secrets/%s", filename)

	knownProjectPath := env.GetVars(ctx).String(EnvVarProjectConfigPath, fallbackProjectConfigPath)
	knownPath := env.GetVars(ctx).String(EnvVarConfigPath, fallbackConfigPath)

	return configutil.OptPaths(
		knownProjectPath,
		knownPath,
	)
}
