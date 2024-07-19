/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
)

func Test_OptAddPaths(t *testing.T) {
	assert := assert.New(t)

	var options ConfigOptions
	assert.Nil(OptAddPaths("foo", "bar")(&options))
	assert.Len(options.FilePaths, 2)
	assert.Equal([]string{"foo", "bar"}, options.FilePaths)

	assert.Nil(OptAddPaths("moo", "loo")(&options))
	assert.Len(options.FilePaths, 4)
	assert.Equal([]string{"foo", "bar", "moo", "loo"}, options.FilePaths)
}

func Test_OptAddPreferredPaths(t *testing.T) {
	assert := assert.New(t)

	var options ConfigOptions
	assert.Nil(OptAddPreferredPaths("foo", "bar")(&options))
	assert.Len(options.FilePaths, 2)
	assert.Equal([]string{"foo", "bar"}, options.FilePaths)

	assert.Nil(OptAddPreferredPaths("moo", "loo")(&options))
	assert.Len(options.FilePaths, 4)
	assert.Equal([]string{"moo", "loo", "foo", "bar"}, options.FilePaths)
}

func Test_OptPaths(t *testing.T) {
	assert := assert.New(t)

	var options ConfigOptions
	assert.Nil(OptPaths("foo", "bar")(&options))
	assert.Len(options.FilePaths, 2)
	assert.Equal([]string{"foo", "bar"}, options.FilePaths)

	assert.Nil(OptPaths("moo", "loo")(&options))
	assert.Len(options.FilePaths, 2)
	assert.Equal([]string{"moo", "loo"}, options.FilePaths)
}

func Test_OptEnv(t *testing.T) {
	assert := assert.New(t)

	var options ConfigOptions
	assert.Empty(options.Env)
	assert.Nil(OptEnv(env.Vars{"FOO": "bar"})(&options))
	assert.NotEmpty(options.Env)
	assert.Len(options.Env, 1)
	assert.Equal("bar", options.Env["FOO"])
}
