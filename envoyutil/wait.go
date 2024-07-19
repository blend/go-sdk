/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package envoyutil

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/retry"
)

// NOTE: Ensure that
//   - `http.Client` satisfies `HTTPGetClient`
//   - `WaitForAdmin.executeOnce` satisfies `retry.Action`
var (
	_ HTTPGetClient      = (*http.Client)(nil)
	_ retry.ActionerFunc = (*WaitForAdmin)(nil).executeOnce
)

var (
	// ErrFailedAttempt is an error class returned when Envoy fails to be
	// ready on a single attempt.
	ErrFailedAttempt = ex.Class("Envoy not yet ready")
	// ErrTimedOut is an error class returned when Envoy fails to be ready
	// after exhausting all attempts.
	ErrTimedOut = ex.Class("Timed out waiting for Envoy to be ready")
)

const (
	// EnvVarWaitFlag is an environment variable which specifies whether
	// a wait function should wait for the Envoy Admin API to be ready.
	EnvVarWaitFlag = "WAIT_FOR_ENVOY"
	// EnvVarAdminPort is an environment variable which provides an override
	// for the Envoy Admin API port.
	EnvVarAdminPort = "ENVOY_ADMIN_PORT"
	// DefaultAdminPort is the default port used for the Envoy Admin API.
	DefaultAdminPort = "15000"
	// EnumStateLive is a `envoy.admin.v3.ServerInfo.State` value indicating
	// the Envoy server is LIVE. Other possible values of this enum are
	// DRAINING, PRE_INITIALIZING and INITIALIZING, but they are not used
	// here.
	// See: https://github.com/envoyproxy/envoy/blob/b867a4dfae32e600ea0a4087dc7925ded5e2ab2a/api/envoy/admin/v3/server_info.proto#L24-L36
	EnumStateLive = "LIVE"
)

// HTTPGetClient captures a small part of the `http.Client` interface needed
// to execute a GET request.
type HTTPGetClient interface {
	Get(url string) (resp *http.Response, err error)
}

// WaitForAdmin encapsulates the settings needed to wait until the Envoy Admin
// API is ready.
type WaitForAdmin struct {
	// Port is the port (on localhost) where the Envoy Admin API is running.
	Port string
	// Sleep is the amount of time to sleep in between failed liveness
	// checks for the Envoy API.
	Sleep time.Duration
	// HTTPClient is the HTTP client to use when sending requests.
	HTTPClient HTTPGetClient
	// Log is an optional logger to be used when executing.
	Log logger.Log
	// Attempt is a counter for the number of attempts that have been made
	// to `executeOnce()`. This makes no attempt at "resetting" or guarding
	// against concurrent usage or re-usage of a `WaitForAdmin` struct.
	Attempt uint32
}

// IsReady makes a single request to the Envoy Admin API and checks if
// the status is ready.
func (wfa *WaitForAdmin) IsReady() bool {
	readyURL := fmt.Sprintf("http://localhost:%s/ready", wfa.Port)
	resp, err := wfa.HTTPClient.Get(readyURL)
	if err != nil {
		logger.MaybeDebugf(wfa.Log, "Envoy is not ready; connection failed: %s", err)
		return false
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.MaybeDebug(wfa.Log, "Envoy is not ready; failed to read response body")
		return false
	}

	if resp.StatusCode != http.StatusOK {
		logger.MaybeDebugf(wfa.Log, "Envoy is not ready; response status code: %d", resp.StatusCode)
		return false
	}

	if string(body) != EnumStateLive+"\n" {
		logger.MaybeDebugf(wfa.Log, "Envoy is not ready; response body: %q", string(body))
		return false
	}

	return true
}

func (wfa *WaitForAdmin) executeOnce(_ context.Context, _ interface{}) (interface{}, error) {
	attempt := atomic.AddUint32(&wfa.Attempt, 1)
	logger.MaybeDebugf(wfa.Log, "Checking if Envoy is ready, attempt %d", attempt)
	if wfa.IsReady() {
		logger.MaybeDebug(wfa.Log, "Envoy is ready")
		return nil, nil
	}

	logger.MaybeDebugf(wfa.Log, "Envoy is not yet ready, sleeping for %s", wfa.Sleep)
	return nil, ErrFailedAttempt
}

// Execute will communicate with the Envoy admin port running on `localhost`,
// which defaults to 15000 but can be overridden with `ENVOY_ADMIN_PORT`. It
// will send `GET /ready` up to 10 times, sleeping for `wfa.Sleep` in between
// if the response is not 200 OK with a body of `LIVE\n`.
func (wfa *WaitForAdmin) Execute(ctx context.Context) error {
	_, err := retry.Retry(
		ctx,
		retry.ActionerFunc(wfa.executeOnce),
		nil,
		retry.OptConstantDelay(wfa.Sleep),
		retry.OptMaxAttempts(10),
	)
	if ex.Is(err, ErrFailedAttempt) {
		return ex.New(ErrTimedOut)
	}
	return err
}

// MaybeWaitForAdmin will check if Envoy is running if the `WAIT_FOR_ENVOY`
// environment variable is set. This will communicate with the Envoy admin
// port running on `localhost`, which defaults to 15000 but can be overridden
// with `ENVOY_ADMIN_PORT`. It will send `GET /ready` up to 10 times, sleeping
// for 1 second in between if the response is not 200 OK with a body of
// `LIVE\n`.
func MaybeWaitForAdmin(log logger.Log) error {
	if !strings.EqualFold(env.Env()[EnvVarWaitFlag], "true") {
		return nil
	}

	hc := &http.Client{Timeout: time.Second}
	wfa := WaitForAdmin{
		Port:       env.Env().String(EnvVarAdminPort, DefaultAdminPort),
		Sleep:      time.Second,
		HTTPClient: hc,
		Log:        log,
		Attempt:    0,
	}

	ctx := context.Background()
	return wfa.Execute(ctx)
}
