/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/blend/go-sdk/r2"
	"github.com/blend/go-sdk/webutil"
)

// Mock sends a mock request to an app.
// It will reset the app Server, Listener, and will set the request host to the listener address
// for a randomized local listener.
func Mock(app *App, req *http.Request, options ...r2.Option) *MockResult {
	var err error
	result := &MockResult{
		App: app,
		Request: &r2.Request{
			Request: req,
		},
	}
	for _, option := range options {
		if err = option(result.Request); err != nil {
			result.Err = err
			return result
		}
	}

	if err := app.StartupTasks(); err != nil {
		result.Err = err
		return result
	}

	if result.Request.Request.URL == nil {
		result.Request.Request.URL = &url.URL{}
	}

	result.Server = httptest.NewUnstartedServer(app)
	result.Server.Config.BaseContext = app.BaseContext
	result.Server.Start()
	result.Request.Closer = result.Close

	parsedServerURL := webutil.MustParseURL(result.Server.URL)
	result.Request.Request.URL.Scheme = parsedServerURL.Scheme
	result.Request.Request.URL.Host = parsedServerURL.Host

	return result
}

// MockMethod sends a mock request with a given method to an app.
// You should use request options to set the body of the request if it's a post or put etc.
func MockMethod(app *App, method, path string, options ...r2.Option) *MockResult {
	req := &http.Request{
		Method: method,
		URL: &url.URL{
			Path: path,
		},
	}
	return Mock(app, req, options...)
}

// MockGet sends a mock get request to an app.
func MockGet(app *App, path string, options ...r2.Option) *MockResult {
	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Path: path,
		},
	}
	return Mock(app, req, options...)
}

// MockPost sends a mock post request to an app.
func MockPost(app *App, path string, body io.ReadCloser, options ...r2.Option) *MockResult {
	req := &http.Request{
		Method: "POST",
		Body:   body,
		URL: &url.URL{
			Path: path,
		},
	}
	return Mock(app, req, options...)
}

// MockPostJSON sends a mock post request with a json body to an app.
func MockPostJSON(app *App, path string, body interface{}, options ...r2.Option) *MockResult {
	contents, _ := json.Marshal(body)
	req := &http.Request{
		Method: "POST",
		Body:   io.NopCloser(bytes.NewReader(contents)),
		URL: &url.URL{
			Path: path,
		},
	}
	return Mock(app, req, options...)
}

// MockResult is a result of a mocked request.
type MockResult struct {
	*r2.Request
	App    *App
	Server *httptest.Server
}

// Close stops the app.
func (mr *MockResult) Close() error {
	mr.Server.Close()
	return nil
}

// MockCtx returns a new mock ctx.
// It is intended to be used in testing.
func MockCtx(method, path string, options ...CtxOption) *Ctx {
	return MockCtxWithBuffer(method, path, new(bytes.Buffer), options...)
}

// MockCtxWithBuffer returns a new mock ctx.
// It is intended to be used in testing.
func MockCtxWithBuffer(method, path string, buf io.Writer, options ...CtxOption) *Ctx {
	return NewCtx(
		webutil.NewMockResponse(buf),
		webutil.NewMockRequest(method, path),
		append(
			[]CtxOption{
				OptCtxApp(new(App)),
				OptCtxDefaultProvider(Text),
			},
			options...,
		)...,
	)
}

// MockSimulateLogin simulates a user login for a given app as mocked request params (i.e. r2 options).
//
// This requires an auth manager to be set on the app.
func MockSimulateLogin(ctx context.Context, app *App, userID string, opts ...r2.Option) []r2.Option {
	sessionID := NewSessionID()
	session := NewSession(userID, sessionID)
	if app.Auth.PersistHandler != nil {
		_ = app.Auth.PersistHandler(ctx, session)
	}
	return append([]r2.Option{
		r2.OptCookieValue(app.Auth.CookieDefaults.Name, sessionID),
	}, opts...)
}
