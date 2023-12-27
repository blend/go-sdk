/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package email

import "context"

// NewMockSender creates a new mock sender.
func NewMockSender() MockSender {
	return MockSender(make(chan Message))
}

// MockSender is a mocked sender.
type MockSender chan Message

// Send sends a mocked message.
func (ms MockSender) Send(ctx context.Context, m Message) error {
	ms <- m
	return nil
}
