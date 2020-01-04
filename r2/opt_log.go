package r2

import (
	"github.com/blend/go-sdk/logger"
)

// OptLog adds OnRequest and OnResponse listeners to log that a call was made.
func OptLog(log logger.Log) Option {
	return func(r *Request) error {
		if err := OptLogRequest(log)(r); err != nil {
			return err
		}
		if err := OptLogResponse(log)(r); err != nil {
			return err
		}
		return nil
	}
}
