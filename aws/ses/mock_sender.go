package ses

import (
	"context"

	"github.com/blend/go-sdk/email"
)

var (
	_ Sender = (*MockSender)(nil)
)

// NewMockSender creates a new mock sender.
func NewMockSender() MockSender {
	return MockSender(make(chan email.Message))
}

// MockSender is a mocked sender.
type MockSender chan email.Message

// Send sends a mocked message.
func (ms MockSender) Send(ctx context.Context, m email.Message) error {
	ms <- m
	return nil
}
