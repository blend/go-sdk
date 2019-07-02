package airbrake

import (
	"context"
	"net/http"

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
	listener := logger.NewErrorEventListener(func(_ context.Context, ee *logger.ErrorEvent) {
		if req, ok := ee.State.(*http.Request); ok {
			client.NotifyWithRequest(ee.Err, req)
		} else if state, ok := ee.State.(Params); ok {
			client.NotifyWithParams(ee.Err, state)
		} else {
			client.Notify(ee.Err)
		}
	})
	log.Listen(logger.Error, ListenerAirbrake, listener)
	log.Listen(logger.Fatal, ListenerAirbrake, listener)
}
