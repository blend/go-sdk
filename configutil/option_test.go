package configutil

import (
	"testing"
	"github.com/blend/go-sdk/assert"
)

func TestOptAddPaths(t *testing.T) {
	assert := assert.New(t)

	var options ConfigOptions
	assert.Nil(OptAddPaths("foo", "bar")(&options))
	assert.Len(options.Paths, 2)
	assert.Equal([]string{"foo", "bar"}, options.Paths)

	assert.Nil(OptAddPaths("moo", "loo")(&options))
	assert.Len(options.Paths, 4)
	assert.Equal([]string{"foo", "bar", "moo", "loo"}, options.Paths)
}

func TestOptAddPreferredPaths(t *testing.T) {
	assert := assert.New(t)

	var options ConfigOptions
	assert.Nil(OptAddPreferredPaths("foo", "bar")(&options))
	assert.Len(options.Paths, 2)
	assert.Equal([]string{"foo", "bar"}, options.Paths)

	assert.Nil(OptAddPreferredPaths("moo", "loo")(&options))
	assert.Len(options.Paths, 4)
	assert.Equal([]string{"moo", "loo", "foo", "bar"}, options.Paths)
}

func TestOptPaths(t *testing.T) {
	assert := assert.New(t)

	var options ConfigOptions
	assert.Nil(OptPaths("foo", "bar")(&options))
	assert.Len(options.Paths, 2)
	assert.Equal([]string{"foo", "bar"}, options.Paths)

	assert.Nil(OptPaths("moo", "loo")(&options))
	assert.Len(options.Paths, 2)
	assert.Equal([]string{"moo", "loo"}, options.Paths)
}

func TestOptResolver(t *testing.T) {
	assert := assert.New(t)

	var options ConfigOptions
	assert.Nil(options.Resolver)
	assert.Nil(OptResolver(func(interface{}) error { return nil })(&options))
	assert.NotNil(options.Resolver)
}
