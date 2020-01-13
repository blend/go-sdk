package sentry

import (
	"context"
	"github.com/blend/go-sdk/logger"
)

// Sender is the type that the sentry client ascribes to.
type Sender interface {
	Notify(context.Context, logger.ErrorEvent)
}
