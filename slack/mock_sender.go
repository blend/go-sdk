/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package slack

import (
	"context"
)

// NewMockWebhookSender creates a new mock sender.
func NewMockWebhookSender() MockWebhookSender {
	return MockWebhookSender(make(chan Message))
}

var (
	_ Sender = (*MockWebhookSender)(nil)
)

// MockWebhookSender is a mocked sender.
type MockWebhookSender chan Message

// Send sends a mocked message.
func (ms MockWebhookSender) Send(ctx context.Context, m Message) error {
	ms <- m
	return nil
}

// SendAndReadResponse sends a mocked message.
func (ms MockWebhookSender) SendAndReadResponse(ctx context.Context, m Message) (*PostMessageResponse, error) {
	ms <- m
	return nil, nil
}

// PostMessage sends a mocked message.
func (ms MockWebhookSender) PostMessage(channel, text string, options ...MessageOption) error {
	m := Message{
		Channel:	channel,
		Text:		text,
	}
	for _, option := range options {
		option(&m)
	}

	ms <- m
	return nil
}

// PostMessageContext sends a mocked message.
func (ms MockWebhookSender) PostMessageContext(ctx context.Context, channel, text string, options ...MessageOption) error {
	m := Message{
		Channel:	channel,
		Text:		text,
	}
	for _, option := range options {
		option(&m)
	}

	ms <- m
	return nil
}
