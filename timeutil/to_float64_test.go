package timeutil

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestToFloat64(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(1550059200000000000, ToFloat64(time.Date(2019, 02, 13, 12, 0, 0, 0, time.UTC)))
}
