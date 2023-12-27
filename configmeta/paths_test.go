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
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/configutil"
	"github.com/blend/go-sdk/env"
)

func Test_Paths_fallbacks(t *testing.T) {
	its := assert.New(t)

	var opts configutil.ConfigOptions

	vars := env.Vars{
		env.VarServiceName: "bar",
	}
	ctx := env.WithVars(context.Background(), vars)
	its.Nil(PathsFileContext(ctx, "foo.yml")(&opts))
	its.Equal([]string{"/var/secrets/projects/bar/foo.yml", "/var/secrets/foo.yml"}, opts.FilePaths)
}

func Test_Paths_fallbacks_project(t *testing.T) {
	its := assert.New(t)

	var opts configutil.ConfigOptions

	vars := env.Vars{
		env.VarServiceName:	"bar",
		env.VarProjectName:	"bar-proj",
	}
	ctx := env.WithVars(context.Background(), vars)
	its.Nil(PathsFileContext(ctx, "foo.yml")(&opts))
	its.Equal([]string{"/var/secrets/projects/bar-proj/foo.yml", "/var/secrets/foo.yml"}, opts.FilePaths)
}

func Test_Paths_env_projectPath(t *testing.T) {
	its := assert.New(t)

	var opts configutil.ConfigOptions

	vars := env.Vars{
		env.VarServiceName:		"bar",
		env.VarProjectName:		"bar-proj",
		EnvVarProjectConfigPath:	"/var/project/secrets/bar/foo.yml",
	}
	ctx := env.WithVars(context.Background(), vars)
	its.Nil(PathsFileContext(ctx, "foo.yml")(&opts))
	its.Equal([]string{"/var/project/secrets/bar/foo.yml", "/var/secrets/foo.yml"}, opts.FilePaths)
}

func Test_Paths_env_configPath(t *testing.T) {
	its := assert.New(t)

	var opts configutil.ConfigOptions

	vars := env.Vars{
		env.VarServiceName:		"bar",
		env.VarProjectName:		"bar-proj",
		EnvVarProjectConfigPath:	"/var/project/secrets/bar/foo.yml",
		EnvVarConfigPath:		"/var/not-secrets/foo.yml",
	}
	ctx := env.WithVars(context.Background(), vars)
	its.Nil(PathsFileContext(ctx, "foo.yml")(&opts))
	its.Equal([]string{"/var/project/secrets/bar/foo.yml", "/var/not-secrets/foo.yml"}, opts.FilePaths)
}
