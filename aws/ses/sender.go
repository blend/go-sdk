package ses

import (
	"context"

	"github.com/blend/go-sdk/email"
)

// Sender is an email sender.
type Sender interface {
	Send(context.Context, email.Message) error
}
