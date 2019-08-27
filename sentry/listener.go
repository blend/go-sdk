package sentry

import (
	"github.com/blend/go-sdk/logger"
)

// AddListeners adds error listeners.
func AddListeners(log logger.Listenable, cfg Config) {
	if log == nil {
		return
	}

	if cfg.DSN != "" {
		client := MustNew(cfg)
		listener := logger.NewErrorEventListener(client.Notify)

		log.Listen(logger.Error, "sentry", listener)
		log.Listen(logger.Fatal, "sentry", listener)
	}
}
