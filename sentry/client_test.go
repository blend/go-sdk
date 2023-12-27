/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package sentry

import (
	"context"
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/webutil"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	c, err := New(Config{
		DSN: "https://example.org",
	})
	assert.NotNil(err)
	assert.Nil(c)

	c, err = New(Config{
		DSN:		"https://foo@example.org/1",
		Environment:	"test",
		ServerName:	"go-sdk-server",
		Release:	"v1.0.0",
		Dist:		"deadbeef",
	})
	assert.Nil(err)
	assert.Equal("test", c.Client.Options().Environment)
	assert.Equal("go-sdk-server", c.Client.Options().ServerName)
	assert.Equal("v1.0.0", c.Client.Options().Release)
	assert.Equal("deadbeef", c.Client.Options().Dist)
}

func TestErrEvent(t *testing.T) {
	assert := assert.New(t)

	event := errEvent(context.Background(), logger.ErrorEvent{
		Flag:	logger.Fatal,
		Err:	ex.New("this is a test", ex.OptMessage("a message")),
		State: &http.Request{
			Method:	"POST",
			Host:	"example.org",
			TLS:	&tls.ConnectionState{},
			URL:	webutil.MustParseURL("https://example.org/foo"),
		},
	})

	assert.NotNil(event)
	assert.NotZero(event.Timestamp)
	assert.Equal(logger.Fatal, event.Level)
	assert.Equal("go", event.Platform)
	assert.Equal(SDK, event.Sdk.Name)
	assert.Equal("this is a test", event.Message)
	assert.NotEmpty(event.Exception)
	assert.NotEmpty(event.Fingerprint)
	assert.NotNil(event.Exception[0].Stacktrace)
	assert.NotEmpty(event.Exception[0].Stacktrace.Frames)
}

func TestErrEvent_Inner(t *testing.T) {
	assert := assert.New(t)

	inner0 := ex.New("this is inner exception 0", ex.OptMessage("inner message 0"))
	inner1 := ex.New("this is inner exception 1", ex.OptMessage("inner message 1"))
	inner2 := ex.New("this is inner exception 2", ex.OptMessage("inner message 2"))
	err := ex.New("this is a test", ex.OptMessage("a message"))

	event := errEvent(context.Background(), logger.ErrorEvent{
		Flag:	logger.Fatal,
		Err:	ex.Nest(err, inner0, inner1, inner2),
		State: &http.Request{
			Method:	"POST",
			Host:	"example.org",
			TLS:	&tls.ConnectionState{},
			URL:	webutil.MustParseURL("https://example.org/foo"),
		},
	})

	assert.NotNil(event)
	assert.NotZero(event.Timestamp)
	assert.Equal(logger.Fatal, event.Level)
	assert.Equal("go", event.Platform)
	assert.Equal(SDK, event.Sdk.Name)
	assert.Equal("this is a test", event.Message)
	assert.Len(event.Exception, 4)
	assert.NotEmpty(event.Fingerprint)
	assert.NotNil(event.Exception[0].Stacktrace)
	assert.NotEmpty(event.Exception[0].Stacktrace.Frames)
}

func TestErrRequest(t *testing.T) {
	assert := assert.New(t)

	res := errRequest(logger.ErrorEvent{})
	assert.Empty(res.URL)

	res = errRequest(logger.ErrorEvent{
		State: &http.Request{
			Method:	"POST",
			Host:	"example.org",
			TLS:	&tls.ConnectionState{},
			URL:	webutil.MustParseURL("https://example.org/foo"),
		},
	})
	assert.Equal("POST", res.Method)
	assert.Equal("https://example.org/foo", res.URL)
}

func TestErrFrames(t *testing.T) {
	assert := assert.New(t)

	err := ex.New("this is only a test")
	assert.NotEmpty(errFrames(err))
}

func TestErrEventFingerprintDefault(t *testing.T) {
	assert := assert.New(t)

	event := errEvent(context.Background(), logger.ErrorEvent{
		Flag:	logger.Fatal,
		Err:	ex.New("this is a test", ex.OptMessage("a message")),
		State: &http.Request{
			Method:	"POST",
			Host:	"example.org",
			TLS:	&tls.ConnectionState{},
			URL:	webutil.MustParseURL("https://example.org/foo"),
		},
	})

	assert.NotNil(event)
	assert.NotEmpty(event.Fingerprint)
	assert.Equal([]string{"this is a test"}, event.Fingerprint)
}

func TestErrEventFingerprintIgnoresPath(t *testing.T) {
	assert := assert.New(t)

	ctx := logger.WithPath(context.Background(), "foo", "bar")

	event := errEvent(ctx, logger.ErrorEvent{
		Flag:	logger.Fatal,
		Err:	ex.New("this is a test", ex.OptMessage("a message")),
		State: &http.Request{
			Method:	"POST",
			Host:	"example.org",
			TLS:	&tls.ConnectionState{},
			URL:	webutil.MustParseURL("https://example.org/foo"),
		},
	})

	assert.NotNil(event)
	assert.NotEmpty(event.Fingerprint)
	assert.Equal([]string{"this is a test"}, event.Fingerprint)
}

func TestErrEventFingerprintOverride(t *testing.T) {
	assert := assert.New(t)

	ctx := WithFingerprint(context.Background(), "test", "fingerprint")

	event := errEvent(ctx, logger.ErrorEvent{
		Flag:	logger.Fatal,
		Err:	ex.New("this is a test", ex.OptMessage("a message")),
		State: &http.Request{
			Method:	"POST",
			Host:	"example.org",
			TLS:	&tls.ConnectionState{},
			URL:	webutil.MustParseURL("https://example.org/foo"),
		},
	})
	assert.NotNil(event)
	assert.NotEmpty(event.Fingerprint)
	assert.Equal([]string{"test", "fingerprint"}, event.Fingerprint)
}
