package cron

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestMinMax(t *testing.T) {
	assert := assert.New(t)

	a := time.Date(2018, 10, 21, 12, 0, 0, 0, time.UTC)
	b := time.Date(2018, 10, 20, 12, 0, 0, 0, time.UTC)

	assert.True(Min(time.Time{}, time.Time{}).IsZero())
	assert.Equal(a, Min(a, time.Time{}))
	assert.Equal(b, Min(time.Time{}, b))
	assert.Equal(b, Min(a, b))
	assert.Equal(b, Min(b, a))

	assert.Equal(a, Max(a, b))
	assert.Equal(a, Max(b, a))
}

func TestConst(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(true, ConstBool(true)())
	assert.Equal(false, ConstBool(false)())

	assert.Equal(123, ConstInt(123)())
	assert.Equal(6*time.Hour, ConstDuration(6*time.Hour)())
	assert.Equal("foo", ConstLabels(map[string]string{
		"bar": "buzz",
		"moo": "foo",
	})()["moo"])
}
