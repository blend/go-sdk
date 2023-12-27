/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

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

// OptLogWithBody adds OnRequest and OnResponse listeners to log that a call was made.
// It will also display the body of the response.
func OptLogWithBody(log logger.Log) Option {
	return func(r *Request) error {
		if err := OptLogRequest(log)(r); err != nil {
			return err
		}
		if err := OptLogResponseWithBody(log)(r); err != nil {
			return err
		}
		return nil
	}
}
