package cron

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestJobInvocationElapsed(t *testing.T) {
	assert := assert.New(t)

	started := time.Now().UTC()

	assert.Equal(200*time.Millisecond, (&JobInvocation{
		Started:  started,
		Complete: started.Add(200 * time.Millisecond),
	}).Elapsed())

	assert.Zero((&JobInvocation{
		Started: started,
	}).Elapsed())
}
