package airbrake

import (
	"context"

	"github.com/blend/go-sdk/logger"
)

const (
	// ListenerAirbrake is the airbrake listener name.
	ListenerAirbrake = "airbrake"
)

// AddListeners adds airbrake listeners.
func AddListeners(log logger.Listenable, cfg Config) {
	if log == nil || cfg.IsZero() {
		return
	}
	client := MustNew(cfg)
	listener := logger.NewErrorEventListener(func(_ context.Context, ee logger.ErrorEvent) {
		if ee.Request != nil {
			client.NotifyWithRequest(ee.Err, ee.Request)
		} else {
			client.Notify(ee.Err)
		}
	})
	log.Listen(logger.Error, ListenerAirbrake, listener)
	log.Listen(logger.Fatal, ListenerAirbrake, listener)
}
