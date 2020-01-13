package cron

import (
	"context"
	"testing"
	"github.com/blend/go-sdk/assert"
)

func Test_JobParametersContext(t *testing.T) {
	assert := assert.New(t)

	final := GetJobParameters(WithJobParameters(context.Background(), JobParameters{
		"foo":  "bar",
		"buzz": "fuzz",
	}))
	assert.Equal("bar", final["foo"])
	assert.Equal("fuzz", final["buzz"])

	assert.Empty(GetJobParameters(context.Background()))
}
