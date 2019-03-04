package timeutil

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestMillis(t *testing.T) {
	assert := assert.New(t)

	duration, err := time.ParseDuration("300ms")
	assert.Nil(err)
	millis := Millis(duration)
	assert.Equal(float64(300), millis)

	duration, err = time.ParseDuration("0ms")
	assert.Nil(err)
	millis = Millis(duration)
	assert.Equal(float64(0), millis)

	duration, err = time.ParseDuration("-2000000ms")
	assert.Nil(err)
	millis = Millis(duration)
	assert.Equal(float64(-2000000), millis)
}
