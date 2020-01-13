package r2

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	"github.com/blend/go-sdk/logger"
)

// OptLogResponse adds an OnResponse listener to log the response of a call.
func OptLogResponse(log logger.Triggerable) Option {
	return OptOnResponse(func(req *http.Request, res *http.Response, started time.Time, err error) error {
		if err != nil {
			return err
		}
		event := NewEvent(FlagResponse,
			OptEventRequest(req),
			OptEventResponse(res),
			OptEventElapsed(time.Now().UTC().Sub(started)),
		)

		log.Trigger(req.Context(), event)
		return nil
	})
}

// OptLogResponseWithBody adds an OnResponse listener to log the response of a call.
// It reads the contents of the response fully before emitting the event.
// Do not use this if the size of the responses can be large.
func OptLogResponseWithBody(log logger.Triggerable) Option {
	return OptOnResponse(func(req *http.Request, res *http.Response, started time.Time, err error) error {
		if err != nil {
			return err
		}
		defer res.Body.Close()

		// read out the buffer in full
		buffer := new(bytes.Buffer)
		if _, err := io.Copy(buffer, res.Body); err != nil {
			return err
		}
		// set the body to the read contents
		res.Body = ioutil.NopCloser(bytes.NewReader(buffer.Bytes()))

		event := NewEvent(FlagResponse,
			OptEventRequest(req),
			OptEventResponse(res),
			OptEventBody(buffer.Bytes()),
			OptEventElapsed(time.Now().UTC().Sub(started)),
		)

		log.Trigger(req.Context(), event)
		return nil
	})
}
