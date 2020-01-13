package airbrake

import (
	"context"
	"net/http"
	"github.com/blend/go-sdk/logger"
)

const (
	// ListenerName is the airbrake listener name.
	ListenerName = "airbrake"
)

// AddListeners adds airbrake listeners.
func AddListeners(log logger.Listenable, cfg Config) {
	if log == nil || cfg.IsZero() {
		return
	}
	client := MustNew(cfg)
	listener := logger.NewErrorEventListener(func(_ context.Context, ee logger.ErrorEvent) {
		if ee.State != nil {
			req, ok := ee.State.(*http.Request)
			if ok {
				client.NotifyWithRequest(ee.Err, req)
				return
			}
		}
		client.Notify(ee.Err)
	})
	log.Listen(logger.Error, ListenerName, listener)
	log.Listen(logger.Fatal, ListenerName, listener)
}
