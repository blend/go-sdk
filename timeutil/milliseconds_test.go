package timeutil

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestMilliseconds(t *testing.T) {
	assert := assert.New(t)
	d := time.Millisecond + time.Microsecond

	assert.Equal(1.001, Milliseconds(d))
}

func TestFromMilliseconds(t *testing.T) {
	assert := assert.New(t)
	expected := time.Millisecond + time.Microsecond
	assert.Equal(expected, FromMilliseconds(1.001))
}
