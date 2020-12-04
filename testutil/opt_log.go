package testutil

import "github.com/blend/go-sdk/logger"

// OptLog sets the suite logger.
func OptLog(log logger.Log) Option {
	return func(s *Suite) {
		s.Log = log
	}
}
