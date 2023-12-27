/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package webutil

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestWebhookSend(t *testing.T) {
	assert := assert.New(t)

	var bodyCorrect, methodCorrect, headerCorrect, contentLengthCorrect bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		bodyCorrect = string(body) == `this is only a test`
		methodCorrect = r.Method == "POST"
		headerCorrect = r.Header.Get("X-Test-Value") == "foo"
		contentLengthCorrect = r.ContentLength == 19

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK!\n")
	}))
	defer ts.Close()

	wh := Webhook{
		URL:	ts.URL,
		Method:	"POST",
		Headers: map[string]string{
			"X-Test-Value": "foo",
		},
		Body:	"this is only a test",
	}

	res, err := wh.Send()
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	assert.True(bodyCorrect)
	assert.True(methodCorrect)
	assert.True(headerCorrect)
	assert.True(contentLengthCorrect)
}
