package env_test

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
)

func TestClear(t *testing.T) {
	assert := assert.New(t)

	vars := env.Vars{
		"Foo": "bar",
	}
	env.SetEnv(vars)
	assert.NotEmpty(env.Env())

	env.Clear()
	assert.Empty(env.Env())
}
