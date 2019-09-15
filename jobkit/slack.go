package jobkit

import (
	"fmt"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/slack"
)

// NewSlackMessage returns a new job started message.
func NewSlackMessage(ji *cron.JobInvocation, options ...slack.MessageOption) slack.Message {
	message := slack.Message{}
	if ji.Err != nil {
		message.Attachments = append(message.Attachments,
			slack.MessageAttachment{
				Text:  fmt.Sprintf("%s %s", ji.JobName, ji.State),
				Color: "#ff0000",
			},
			slack.MessageAttachment{
				Text:  fmt.Sprintf("error: %+v", ji.Err),
				Color: "#ff0000",
			},
		)
	} else {
		message.Attachments = append(message.Attachments,
			slack.MessageAttachment{
				Text:  fmt.Sprintf("%s %s", ji.JobName, ji.State),
				Color: "#00ff00",
			},
		)
	}

	if ji.Elapsed > 0 {
		message.Attachments = append(message.Attachments,
			slack.MessageAttachment{
				Text: fmt.Sprintf("%v elapsed", ji.Elapsed),
			},
		)
	}

	return slack.ApplyMessageOptions(message, options...)
}
