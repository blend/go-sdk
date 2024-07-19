/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/blend/go-sdk/logger"
)

// OptLogResponse adds an OnResponse listener to log the response of a call.
func OptLogResponse(log logger.Triggerable) Option {
	return OptOnResponse(func(req *http.Request, res *http.Response, startedUTC time.Time, err error) error {
		if err != nil {
			return nil
		}
		event := NewEvent(FlagResponse,
			OptEventRequest(req),
			OptEventResponse(res),
			OptEventElapsed(time.Now().UTC().Sub(startedUTC)),
		)

		logger.MaybeTriggerContext(req.Context(), log, event)
		return nil
	})
}

// OptLogResponseWithBody adds an OnResponse listener to log the response of a call.
// It reads the contents of the response fully before emitting the event.
// Do not use this if the size of the responses can be large.
func OptLogResponseWithBody(log logger.Triggerable) Option {
	return OptOnResponse(func(req *http.Request, res *http.Response, started time.Time, err error) error {
		if err != nil {
			return nil
		}
		defer res.Body.Close()

		// read out the buffer in full
		buffer := new(bytes.Buffer)
		if _, err := io.Copy(buffer, res.Body); err != nil {
			return err
		}
		// set the body to the read contents
		res.Body = io.NopCloser(bytes.NewReader(buffer.Bytes()))

		event := NewEvent(FlagResponse,
			OptEventRequest(req),
			OptEventResponse(res),
			OptEventBody(buffer.Bytes()),
			OptEventElapsed(time.Now().UTC().Sub(started)),
		)

		logger.MaybeTriggerContext(req.Context(), log, event)
		return nil
	})
}
