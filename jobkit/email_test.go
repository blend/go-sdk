package jobkit

import (
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/email"
)

func TestNewEmailMessage(t *testing.T) {
	assert := assert.New(t)

	message, err := NewEmailMessage(email.Message{}, &cron.JobInvocation{
		JobName: "test",
		Elapsed: time.Millisecond,
		State:   cron.JobInvocationStateComplete,
		Output: &cron.OutputBuffer{
			Chunks: []cron.OutputChunk{
				{Data: []byte("this is a test")},
				{Data: []byte("this is another test")},
			},
		},
	},
		email.OptFrom("jobkit@blend.com"),
		email.OptTo("foo@bar.com"),
		email.OptCC("baileydog@blend.com"),
	)
	assert.Nil(err)
	assert.Equal("test :: complete", message.Subject)
	assert.NotEmpty(message.From)
	assert.Equal("jobkit@blend.com", message.From)
	assert.NotEmpty(message.To)
	assert.Equal("foo@bar.com", message.To[0])
	assert.NotEmpty(message.CC)
	assert.Equal("baileydog@blend.com", message.CC[0])
	assert.NotEmpty(message.HTMLBody)
	assert.Contains(message.HTMLBody, "this is a test")
	assert.Contains(message.HTMLBody, "this is another test")
	assert.NotEmpty(message.TextBody)
	assert.Contains(message.TextBody, "this is a test")
	assert.Contains(message.TextBody, "this is another test")
}
