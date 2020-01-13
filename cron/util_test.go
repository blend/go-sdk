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
