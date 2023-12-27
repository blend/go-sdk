/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package r2

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

// NOTE: Ensure that mockLogger satisfies the logger.Triggerable interface.
var (
	_ logger.Triggerable = (*mockLogger)(nil)
)

type mockLogger struct {
	Events []logger.Event
}

func (ml *mockLogger) TriggerContext(ctx context.Context, e logger.Event) {
	ml.Events = append(ml.Events, e)
}

func TestOptLogResponse(t *testing.T) {
	assert := assert.New(t)
	ml := &mockLogger{}
	opt := OptLogResponse(ml)
	e := logResponseHelper(assert, ml, opt, "OK!\n")
	assert.Equal(e, Event{
		Flag:     FlagResponse,
		Request:  e.Request,
		Response: e.Response,
		Body:     nil,
		Elapsed:  e.Elapsed,
	})
}

func TestOptLogResponseWithBody(t *testing.T) {
	assert := assert.New(t)
	ml := &mockLogger{}
	opt := OptLogResponseWithBody(ml)
	body := "This is the response body\n"
	e := logResponseHelper(assert, ml, opt, body)
	assert.Equal(e, Event{
		Flag:     FlagResponse,
		Request:  e.Request,
		Response: e.Response,
		Body:     []byte(body),
		Elapsed:  e.Elapsed,
	})
	assert.Equal(e.Request.ContentLength, 0)
	assert.Equal(e.Response.ContentLength, 26)
}

func logResponseHelper(a *assert.Assertions, ml *mockLogger, opt Option, body string) Event {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, body)
	}))
	defer server.Close()

	r := New(server.URL, opt)
	res, err := r.Do()
	a.Nil(err)
	a.NotNil(res)
	a.Equal(res.StatusCode, http.StatusOK)
	bodyBytes, err := io.ReadAll(res.Body)
	a.Nil(err)
	a.Equal(bodyBytes, []byte(body))

	// Make sure the event was triggered.
	a.Len(ml.Events, 1)
	e, ok := ml.Events[0].(Event)
	a.True(ok)
	return e
}
