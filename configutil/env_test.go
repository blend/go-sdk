/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import (
	"context"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/uuid"
)

func createEnvVarsContext(key, value string) context.Context {
	options := ConfigOptions{
		Env: env.Vars{key: value},
	}
	return options.Background()
}

func emptyEnvVarsContext() context.Context {
	return ConfigOptions{}.Background()
}

func TestEnv(t *testing.T) {
	assert := assert.New(t)

	key := uuid.V4().String()
	ctx := emptyEnvVarsContext()

	stringValue, err := Env(key).String(ctx)
	assert.Nil(err)
	assert.Nil(stringValue)

	ctx = createEnvVarsContext(key, "foo")
	assert.NotNil(env.GetVars(ctx))
	assert.NotEmpty(env.GetVars(ctx))
	assert.Equal("foo", env.GetVars(ctx).String(key))

	stringValue, err = Env(key).String(ctx)
	assert.Nil(err)
	assert.NotNil(stringValue)
	assert.Equal("foo", *stringValue)

	ctx = emptyEnvVarsContext()
	stringsValue, err := Env(key).Strings(ctx)
	assert.Nil(err)
	assert.Nil(stringsValue)

	ctx = createEnvVarsContext(key, "foo,bar")
	stringsValue, err = Env(key).Strings(ctx)
	assert.Nil(err)
	assert.NotEmpty(stringsValue)
	assert.Equal([]string{"foo", "bar"}, stringsValue)

	ctx = emptyEnvVarsContext()
	boolValue, err := Env(key).Bool(ctx)
	assert.Nil(err)
	assert.Nil(boolValue)

	ctx = createEnvVarsContext(key, "true")
	boolValue, err = Env(key).Bool(ctx)
	assert.Nil(err)
	assert.NotNil(boolValue)
	assert.Equal(true, *boolValue)

	ctx = emptyEnvVarsContext()
	intValue, err := Env(key).Int(ctx)
	assert.Nil(err)
	assert.Nil(intValue)

	ctx = createEnvVarsContext(key, "bad value")
	intValue, err = Env(key).Int(ctx)
	assert.NotNil(err)
	assert.Nil(intValue)

	ctx = createEnvVarsContext(key, "4321")
	intValue, err = Env(key).Int(ctx)
	assert.Nil(err)
	assert.NotNil(intValue)
	assert.Equal(4321, *intValue)

	ctx = emptyEnvVarsContext()
	floatValue, err := Env(key).Float64(ctx)
	assert.Nil(err)
	assert.Nil(floatValue)

	ctx = createEnvVarsContext(key, "bad value")
	floatValue, err = Env(key).Float64(ctx)
	assert.NotNil(err)
	assert.Nil(floatValue)

	ctx = createEnvVarsContext(key, "4321")
	floatValue, err = Env(key).Float64(ctx)
	assert.Nil(err)
	assert.NotNil(floatValue)
	assert.Equal(4321, *floatValue)

	ctx = emptyEnvVarsContext()
	durationValue, err := Env(key).Duration(ctx)
	assert.Nil(err)
	assert.Nil(durationValue)

	ctx = createEnvVarsContext(key, "bad value")
	durationValue, err = Env(key).Duration(ctx)
	assert.NotNil(err)
	assert.Nil(durationValue)

	ctx = createEnvVarsContext(key, "10s")
	durationValue, err = Env(key).Duration(ctx)
	assert.Nil(err)
	assert.NotNil(durationValue)
	assert.Equal(10*time.Second, *durationValue)
}
