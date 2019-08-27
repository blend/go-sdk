package sentry

import (
	"github.com/blend/go-sdk/logger"
)

const (
	// ListenerName is the sentry listener name.
	ListenerName = "airbrake"
)

// AddListeners adds error listeners.
func AddListeners(log logger.Listenable, cfg Config) {
	if log == nil || cfg.IsZero() {
		return
	}
	client := MustNew(cfg)
	listener := logger.NewErrorEventListener(client.Notify)
	log.Listen(logger.Error, ListenerName, listener)
	log.Listen(logger.Fatal, ListenerName, listener)
}
