package r2

import (
	"net/http"

	"github.com/blend/go-sdk/logger"
)

// OptLogRequest adds OnRequest and OnResponse listeners to log that a call was made.
func OptLogRequest(log logger.Log) Option {
	return OptOnRequest(func(req *http.Request) error {
		logger.MaybeTriggerContext(req.Context(), log, NewEvent(Flag, OptEventRequest(req)))
		return nil
	})
}
