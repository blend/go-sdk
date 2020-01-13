package r2

import (
	"net/http"
	"github.com/blend/go-sdk/logger"
)

const (
	// maxLogBytes is the maximum number of bytes to log from a response.
	// it is currently set to 1mb.
	maxLogBytes = 1 << 20
)

// OptLog adds an OnRequest listener to log that a call was made.
func OptLog(log logger.Log) Option {
	return OptOnRequest(func(req *http.Request) error {
		event := NewEvent(Flag,
			OptEventRequest(req))
		log.Trigger(req.Context(), event)
		return nil
	})
}
