package cron

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestJobInvocationElapsed(t *testing.T) {
	assert := assert.New(t)

	started := time.Now().UTC()

	assert.Equal(200*time.Millisecond, JobInvocation{
		Started:   started,
		Errored:   started.Add(500 * time.Millisecond),
		Timeout:   started.Add(400 * time.Millisecond),
		Cancelled: started.Add(300 * time.Millisecond),
		Complete:  started.Add(200 * time.Millisecond),
	}.Elapsed())

	assert.Equal(300*time.Millisecond, JobInvocation{
		Started:   started,
		Errored:   started.Add(500 * time.Millisecond),
		Timeout:   started.Add(400 * time.Millisecond),
		Cancelled: started.Add(300 * time.Millisecond),
	}.Elapsed())

	assert.Equal(400*time.Millisecond, JobInvocation{
		Started: started,
		Errored: started.Add(500 * time.Millisecond),
		Timeout: started.Add(400 * time.Millisecond),
	}.Elapsed())

	assert.Equal(500*time.Millisecond, JobInvocation{
		Started: started,
		Errored: started.Add(500 * time.Millisecond),
	}.Elapsed())
}
