package webutil

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestWebhookIsZero(t *testing.T) {
	assert := assert.New(t)

	wh := Webhook{}
	assert.True(wh.IsZero())

	wh = Webhook{URL: "http://localhost:12345/foo"}
	assert.False(wh.IsZero())
}

func TestWebhookMethodOrDefault(t *testing.T) {
	assert := assert.New(t)

	wh := Webhook{}
	assert.Equal("GET", wh.MethodOrDefault())

	wh = Webhook{Method: "POST"}
	assert.Equal("POST", wh.MethodOrDefault())
}

func TestWebhookSend(t *testing.T) {
	assert := assert.New(t)

	// Invalid URL
	wh := Webhook{URL: "\r"}
	res, err := wh.Send()
	assert.Nil(res)
	assert.NotNil(err)
	assert.Equal("parse \r: net/url: invalid control character in URL", err.Error())

	var bodyCorrect, methodCorrect, headerCorrect, contentLengthCorrect bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
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

	wh = Webhook{
		URL:    ts.URL,
		Method: "POST",
		Headers: map[string]string{
			"X-Test-Value": "foo",
		},
		Body: "this is only a test",
	}

	res, err = wh.Send()
	assert.Nil(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	assert.True(bodyCorrect)
	assert.True(methodCorrect)
	assert.True(headerCorrect)
	assert.True(contentLengthCorrect)

	// Unreachable URL
	wh = Webhook{URL: "http://localhost:11001/foo"}
	res, err = wh.Send()
	assert.Nil(res)
	assert.NotNil(err)
	assert.Equal("Get http://localhost:11001/foo: dial tcp 127.0.0.1:11001: connect: connection refused", err.Error())
}
