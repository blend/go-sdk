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

	stringValue, err := Env(ctx, key).String()
	assert.Nil(err)
	assert.Nil(stringValue)

	ctx = createEnvVarsContext(key, "foo")
	assert.NotNil(GetEnvVars(ctx))
	assert.NotEmpty(GetEnvVars(ctx))
	assert.Equal("foo", GetEnvVars(ctx).String(key))

	stringValue, err = Env(ctx, key).String()
	assert.Nil(err)
	assert.NotNil(stringValue)
	assert.Equal("foo", *stringValue)

	ctx = emptyEnvVarsContext()
	stringsValue, err := Env(ctx, key).Strings()
	assert.Nil(err)
	assert.Nil(stringsValue)

	ctx = createEnvVarsContext(key, "foo,bar")
	stringsValue, err = Env(ctx, key).Strings()
	assert.Nil(err)
	assert.NotEmpty(stringsValue)
	assert.Equal([]string{"foo", "bar"}, stringsValue)

	ctx = emptyEnvVarsContext()
	boolValue, err := Env(ctx, key).Bool()
	assert.Nil(err)
	assert.Nil(boolValue)

	ctx = createEnvVarsContext(key, "true")
	boolValue, err = Env(ctx, key).Bool()
	assert.Nil(err)
	assert.NotNil(boolValue)
	assert.Equal(true, *boolValue)

	ctx = emptyEnvVarsContext()
	intValue, err := Env(ctx, key).Int()
	assert.Nil(err)
	assert.Nil(intValue)

	ctx = createEnvVarsContext(key, "bad value")
	intValue, err = Env(ctx, key).Int()
	assert.NotNil(err)
	assert.Nil(intValue)

	ctx = createEnvVarsContext(key, "4321")
	intValue, err = Env(ctx, key).Int()
	assert.Nil(err)
	assert.NotNil(intValue)
	assert.Equal(4321, *intValue)

	ctx = emptyEnvVarsContext()
	floatValue, err := Env(ctx, key).Float64()
	assert.Nil(err)
	assert.Nil(floatValue)

	ctx = createEnvVarsContext(key, "bad value")
	floatValue, err = Env(ctx, key).Float64()
	assert.NotNil(err)
	assert.Nil(floatValue)

	ctx = createEnvVarsContext(key, "4321")
	floatValue, err = Env(ctx, key).Float64()
	assert.Nil(err)
	assert.NotNil(floatValue)
	assert.Equal(4321, *floatValue)

	ctx = emptyEnvVarsContext()
	durationValue, err := Env(ctx, key).Duration()
	assert.Nil(err)
	assert.Nil(durationValue)

	ctx = createEnvVarsContext(key, "bad value")
	durationValue, err = Env(ctx, key).Duration()
	assert.NotNil(err)
	assert.Nil(durationValue)

	ctx = createEnvVarsContext(key, "10s")
	durationValue, err = Env(ctx, key).Duration()
	assert.Nil(err)
	assert.NotNil(durationValue)
	assert.Equal(10*time.Second, *durationValue)
}
