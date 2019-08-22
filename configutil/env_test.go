package configutil

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/uuid"
)

func TestEnv(t *testing.T) {
	assert := assert.New(t)

	defer env.Restore()

	key := uuid.V4().String()
	env.Env().Delete(key)

	stringValue, err := Env(key).String()
	assert.Nil(err)
	assert.Nil(stringValue)

	env.Env().Set(key, "foo")
	stringValue, err = Env(key).String()
	assert.Nil(err)
	assert.NotNil(stringValue)
	assert.Equal("foo", *stringValue)

	env.Env().Delete(key)

	stringsValue, err := Env(key).Strings()
	assert.Nil(err)
	assert.Nil(stringsValue)

	env.Env().Set(key, "foo,bar")
	stringsValue, err = Env(key).Strings()
	assert.Nil(err)
	assert.NotEmpty(stringsValue)
	assert.Equal([]string{"foo", "bar"}, stringsValue)

	env.Env().Delete(key)

	boolValue, err := Env(key).Bool()
	assert.Nil(err)
	assert.Nil(boolValue)

	env.Env().Set(key, "true")
	boolValue, err = Env(key).Bool()
	assert.Nil(err)
	assert.NotNil(boolValue)
	assert.Equal(true, *boolValue)

	env.Env().Delete(key)

	intValue, err := Env(key).Int()
	assert.Nil(err)
	assert.Nil(intValue)

	env.Env().Set(key, "bad value")
	intValue, err = Env(key).Int()
	assert.NotNil(err)
	assert.Nil(intValue)

	env.Env().Set(key, "4321")
	intValue, err = Env(key).Int()
	assert.Nil(err)
	assert.NotNil(intValue)
	assert.Equal(4321, *intValue)

	env.Env().Delete(key)

	floatValue, err := Env(key).Float64()
	assert.Nil(err)
	assert.Nil(floatValue)

	env.Env().Set(key, "bad value")
	floatValue, err = Env(key).Float64()
	assert.NotNil(err)
	assert.Nil(floatValue)

	env.Env().Set(key, "4321")
	floatValue, err = Env(key).Float64()
	assert.Nil(err)
	assert.NotNil(floatValue)
	assert.Equal(4321, *floatValue)

	env.Env().Delete(key)

	durationValue, err := Env(key).Duration()
	assert.Nil(err)
	assert.Nil(durationValue)

	env.Env().Set(key, "bad value")
	durationValue, err = Env(key).Duration()
	assert.NotNil(err)
	assert.Nil(durationValue)

	env.Env().Set(key, "10s")
	durationValue, err = Env(key).Duration()
	assert.Nil(err)
	assert.NotNil(durationValue)
	assert.Equal(10*time.Second, *durationValue)
}
