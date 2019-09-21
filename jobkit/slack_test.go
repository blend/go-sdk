package jobkit

import (
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/slack"
	"github.com/blend/go-sdk/uuid"
)

func TestNewSlackMessage(t *testing.T) {
	assert := assert.New(t)

	id := uuid.V4().String()
	jobName := uuid.V4().String()
	message := NewSlackMessage(cron.FlagComplete, &cron.JobInvocation{
		ID:      id,
		JobName: jobName,
		State:   cron.JobInvocationStateComplete,
		Elapsed: time.Second,
		Err:     fmt.Errorf("this is just a test"),
	}, func(m *slack.Message) { m.AsUser = true })
	assert.True(message.AsUser)
	assert.NotEmpty(message.Attachments)
	assert.Contains(message.Attachments[0].Text, jobName+" "+cron.FlagComplete)
	assert.Contains(message.Attachments[1].Text, "this is just a test")
	assert.Contains(message.Attachments[2].Text, "1s elapsed")
}
