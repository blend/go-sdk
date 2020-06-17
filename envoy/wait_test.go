package envoy_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/envoy"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/web"
)

// NOTE: Ensure that
//       - `TimeoutError` satisfies `error`
//       - `BadReadCloser` satisfies `io.ReadCloser`
//       - `MockHTTPGetClient` satisfies `envoy.HTTPGetClient`
var (
	_ error               = (*TimeoutError)(nil)
	_ io.ReadCloser       = (*BadReadCloser)(nil)
	_ envoy.HTTPGetClient = (*MockHTTPGetClient)(nil)
)

func TestMaybeWaitForAdmin(t *testing.T) {
	it := assert.New(t)

	defer env.Restore()
	env.SetEnv(env.New())

	// No-op (WAIT_FOR_ENVOY is not set.)
	var logBuffer bytes.Buffer
	log := InMemoryLog(&logBuffer)
	err := envoy.MaybeWaitForAdmin(log)
	it.Nil(err)
	it.Empty(logBuffer.Bytes())
	logBuffer.Reset()

	// Happy-path; WAIT_FOR_ENVOY / ENVOY_ADMIN_PORT set.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, envoy.EnumStateLive+"\n")
	}))
	defer server.Close()

	port := strings.TrimPrefix(server.URL, "http://127.0.0.1:")
	env.Env().Set(envoy.EnvVarWaitFlag, "true")
	env.Env().Set(envoy.EnvVarAdminPort, port)
	err = envoy.MaybeWaitForAdmin(log)
	it.Nil(err)
	expected := strings.Join([]string{
		"[debug] Checking if Envoy is ready, attempt 1",
		"[debug] Envoy is ready",
		"",
	}, "\n")
	it.Equal(expected, string(logBuffer.Bytes()))
	logBuffer.Reset()
}

func TestWaitForAdminExecute(t *testing.T) {
	it := assert.New(t)

	// Failure with error that isn't timeout or connection error.
	mhgc := &MockHTTPGetClient{Error: ex.New("known failure")}
	wfa := envoy.WaitForAdmin{HTTPClient: mhgc}
	err := wfa.Execute(context.TODO())
	it.True(ex.Is(err, envoy.ErrTimedOut))

	// Repeated failures with timeout
	ue := &url.Error{
		Op:  "Get",
		URL: "http://localhost:15000/ready",
		Err: &TimeoutError{},
	}
	mhgc = &MockHTTPGetClient{Error: ue}
	wfa = envoy.WaitForAdmin{HTTPClient: mhgc, Sleep: time.Nanosecond}
	err = wfa.Execute(context.TODO())
	it.True(ex.Is(err, envoy.ErrTimedOut))

	// Success after repeated failures.
	var logBuffer bytes.Buffer
	log := InMemoryLog(&logBuffer)
	mhgc = &MockHTTPGetClient{
		Error:       ue,
		SwitchAfter: 3,
		SwitchResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(envoy.EnumStateLive + "\n"))),
		},
	}
	wfa = envoy.WaitForAdmin{Log: log, HTTPClient: mhgc, Sleep: time.Nanosecond}
	err = wfa.Execute(context.TODO())
	it.Nil(err)
	expected := strings.Join([]string{
		"[debug] Checking if Envoy is ready, attempt 1",
		timeoutError(),
		"[debug] Envoy is not yet ready, sleeping for 1ns",
		"[debug] Checking if Envoy is ready, attempt 2",
		timeoutError(),
		"[debug] Envoy is not yet ready, sleeping for 1ns",
		"[debug] Checking if Envoy is ready, attempt 3",
		"[debug] Envoy is ready",
		"",
	}, "\n")
	it.Equal(expected, string(logBuffer.Bytes()))
}

// getMinorVersion attempts to extract the minor (`X`) version in a version
// string `go1.X` or `go1.X.Y`, e.g. `go1.14.4`. If the version is not of
// this form, a negative number will be returned.
func getMinorVersion() int {
	version := runtime.Version()
	if !strings.HasPrefix(version, "go1.") {
		return -1
	}

	parts := strings.SplitN(version, ".", 3)
	if len(parts) < 2 {
		return -2
	}

	mv, err := strconv.Atoi(parts[1])
	if err != nil {
		return -3
	}

	return mv
}

func timeoutError() string {
	mv := getMinorVersion()
	if mv < 14 {
		return `[debug] Envoy is not ready; connection failed: Get http://localhost:15000/ready: TimeoutError`
	}
	return `[debug] Envoy is not ready; connection failed: Get "http://localhost:15000/ready": TimeoutError`
}

func TestIsReady(t *testing.T) {
	it := assert.New(t)

	responses := make(chan web.RawResult, 1)
	// Happy-path; WAIT_FOR_ENVOY / ENVOY_ADMIN_PORT set.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result := <-responses
		w.WriteHeader(result.StatusCode)
		_, _ = w.Write(result.Response)
	}))
	defer server.Close()

	port := strings.TrimPrefix(server.URL, "http://127.0.0.1:")
	wfa := envoy.WaitForAdmin{
		Port:       port,
		Sleep:      time.Nanosecond,
		HTTPClient: &http.Client{Timeout: time.Second},
	}

	// Non-200 response code.
	responses <- web.RawResult{
		Response:   []byte("PRE_INITIALIZING\n"),
		StatusCode: http.StatusServiceUnavailable,
	}
	ok := wfa.IsReady()
	it.False(ok)

	// 200 response code, but invalid body
	responses <- web.RawResult{
		Response:   []byte("INITIALIZING\n"),
		StatusCode: http.StatusOK,
	}
	ok = wfa.IsReady()
	it.False(ok)

	// Error reading response body.
	bodyErr := ex.New("Filesystem oops")
	body := &BadReadCloser{Error: bodyErr}
	mhgc := &MockHTTPGetClient{Response: &http.Response{Body: body}}
	wfa = envoy.WaitForAdmin{
		Port:       port,
		Sleep:      time.Nanosecond,
		HTTPClient: mhgc,
	}
	ok = wfa.IsReady()
	it.False(ok)
}

type MockHTTPGetClient struct {
	Response *http.Response
	Error    error
	// CallCount tracks the number of times `Get()` has been called.
	CallCount uint32

	// SwitchAfter is a `CallCount` target. Once the `CallCount` reaches this
	// value, the mocked response from `Get()` will change from `Response, Error`
	// to `SwitchResponse, SwitchError`.
	SwitchAfter    uint32
	SwitchResponse *http.Response
	SwitchError    error
}

func (mhgc *MockHTTPGetClient) Get(url string) (resp *http.Response, err error) {
	count := atomic.AddUint32(&mhgc.CallCount, 1)
	if mhgc.SwitchAfter > 0 && count >= mhgc.SwitchAfter {
		return mhgc.SwitchResponse, mhgc.SwitchError
	}

	return mhgc.Response, mhgc.Error
}

type TimeoutError struct {
}

func (te TimeoutError) Timeout() bool {
	return true
}

func (te TimeoutError) Error() string {
	return "TimeoutError"
}

type BadReadCloser struct {
	Error error
}

func (brc *BadReadCloser) Read(p []byte) (n int, err error) {
	return 0, brc.Error
}

func (brc *BadReadCloser) Close() error {
	return brc.Error
}
