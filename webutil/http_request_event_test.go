/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package webutil

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func TestNewHTTPRequestEvent(t *testing.T) {
	assert := assert.New(t)

	hre := NewHTTPRequestEvent(nil,
		OptHTTPRequestRequest(&http.Request{Method: "foo", URL: &url.URL{Scheme: "https", Host: "localhost", Path: "/foo/example-string"}}),
		OptHTTPRequestContentEncoding("utf-8"),
		OptHTTPRequestContentLength(1337),
		OptHTTPRequestContentType("text/html"),
		OptHTTPRequestElapsed(time.Second),
		OptHTTPRequestRoute("/foo/:bar"),
		OptHTTPRequestStatusCode(http.StatusOK),
		OptHTTPRequestHeader(http.Header{"X-Bad": []string{"nope", "definitely nope"}}),
	)

	assert.Equal("foo", hre.Request.Method)
	assert.Equal("utf-8", hre.ContentEncoding)
	assert.Equal(1337, hre.ContentLength)
	assert.Equal("text/html", hre.ContentType)
	assert.Equal(time.Second, hre.Elapsed)
	assert.Equal("/foo/:bar", hre.Route)
	assert.Equal(http.StatusOK, hre.StatusCode)
	assert.Equal("nope", hre.Header.Get("X-Bad"))

	noColor := logger.NewTextOutputFormatter(logger.OptTextNoColor())
	buf := new(bytes.Buffer)
	hre.WriteText(noColor, buf)
	assert.NotContains(buf.String(), "/foo/:bar")
	assert.Contains(buf.String(), "/foo/example-string")
	assert.NotContains(buf.String(), "X-Bad", "response headers should not be written to text output")
	assert.NotContains(buf.String(), "definitely nope", "response headers should not be written to text output")

	contents, err := json.Marshal(hre.Decompose())
	assert.Nil(err)
	assert.Contains(string(contents), "/foo/:bar")

	assert.NotContains(string(contents), "X-Bad", "response headers should not be written to json output")
	assert.NotContains(string(contents), "definitely nope", "response headers should not be written to json output")
}

func TestHTTPRequestEventListener(t *testing.T) {
	assert := assert.New(t)

	var didCall bool
	listener := NewHTTPRequestEventListener(func(_ context.Context, hre HTTPRequestEvent) {
		didCall = true
	})
	listener(context.Background(), logger.NewMessageEvent(logger.Info, "test"))
	assert.False(didCall)
	listener(context.Background(), NewHTTPRequestEvent(nil))
	assert.True(didCall)
}
