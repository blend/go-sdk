package configutil

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
)

func TestOptAddPaths(t *testing.T) {
	assert := assert.New(t)

	var options ConfigOptions
	assert.Nil(OptAddFilePaths("foo", "bar")(&options))
	assert.Len(options.FilePaths, 2)
	assert.Equal([]string{"foo", "bar"}, options.FilePaths)

	assert.Nil(OptAddFilePaths("moo", "loo")(&options))
	assert.Len(options.FilePaths, 4)
	assert.Equal([]string{"foo", "bar", "moo", "loo"}, options.FilePaths)
}

func TestOptAddPreferredPaths(t *testing.T) {
	assert := assert.New(t)

	var options ConfigOptions
	assert.Nil(OptAddPreferredFilePaths("foo", "bar")(&options))
	assert.Len(options.FilePaths, 2)
	assert.Equal([]string{"foo", "bar"}, options.FilePaths)

	assert.Nil(OptAddPreferredFilePaths("moo", "loo")(&options))
	assert.Len(options.FilePaths, 4)
	assert.Equal([]string{"moo", "loo", "foo", "bar"}, options.FilePaths)
}

func TestOptPaths(t *testing.T) {
	assert := assert.New(t)

	var options ConfigOptions
	assert.Nil(OptFilePaths("foo", "bar")(&options))
	assert.Len(options.FilePaths, 2)
	assert.Equal([]string{"foo", "bar"}, options.FilePaths)

	assert.Nil(OptFilePaths("moo", "loo")(&options))
	assert.Len(options.FilePaths, 2)
	assert.Equal([]string{"moo", "loo"}, options.FilePaths)
}

func TestOptEnv(t *testing.T) {
	assert := assert.New(t)

	var options ConfigOptions
	assert.Empty(options.Env)
	assert.Nil(OptEnv(env.Vars{"FOO": "bar"})(&options))
	assert.NotEmpty(options.Env)
	assert.Len(options.Env, 1)
	assert.Equal("bar", options.Env["FOO"])
}
