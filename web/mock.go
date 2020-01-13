package web

import (
	"bytes"
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
			Request: *req,
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

	if result.Request.URL == nil {
		result.Request.URL = &url.URL{}
	}

	result.Server = httptest.NewServer(app)

	parsedServerURL := webutil.MustParseURL(result.Server.URL)
	result.Request.URL.Scheme = parsedServerURL.Scheme
	result.Request.URL.Host = parsedServerURL.Host

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
	return NewCtx(webutil.NewMockResponse(new(bytes.Buffer)), webutil.NewMockRequest(method, path), options...)
}
