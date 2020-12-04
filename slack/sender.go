package slack

import "context"

// Sender is a type that can send slack messages.
type Sender interface {
	Send(ctx context.Context, msg Message) error
	SendAndReadResponse(ctx context.Context, msg Message) (*PostMessageResponse, error)
	PostMessage(channel string, messageText string, opts ...MessageOption) error
	PostMessageContext(ctx context.Context, channel string, messageText string, opts ...MessageOption) error
}
