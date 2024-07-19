/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package status

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/configmeta"
	"github.com/blend/go-sdk/web"
)

func Test_Controller_getStatus(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	statusController := NewController(
		OptConfig(configmeta.Meta{
			ServiceEnv:  "test",
			ServiceName: "test-service",
			Version:     "1.2.3",
		}),
	)
	app := web.MustNew()
	app.Register(statusController)

	var res configmeta.Meta
	meta, err := web.MockGet(app, "/status").JSON(&res)
	its.Nil(err)
	its.Equal(http.StatusOK, meta.StatusCode)

	its.Equal("test", res.ServiceEnv)
	its.Equal("test-service", res.ServiceName)
	its.Equal("1.2.3", res.Version)
}

func Test_Controller_getStatusSLA(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	var fooCalls, barCalls int32
	statusController := NewController(
		OptCheck("foo", CheckerFunc(func(_ context.Context) error {
			atomic.AddInt32(&fooCalls, 1)
			return nil
		})),
		OptCheck("bar", CheckerFunc(func(_ context.Context) error {
			atomic.AddInt32(&barCalls, 1)
			return nil
		})),
	)

	app := web.MustNew()
	app.Register(statusController)

	var slaOutput FreeformResult
	meta, err := web.MockGet(app, "/status/sla").JSON(&slaOutput)
	its.Nil(err)
	its.Equal(http.StatusOK, meta.StatusCode)
	its.True(slaOutput["foo"])
	its.True(slaOutput["bar"])

	meta, err = web.MockGet(app, "/status/sla").JSON(&slaOutput)
	its.Nil(err)
	its.Equal(http.StatusOK, meta.StatusCode)
	its.True(slaOutput["foo"])
	its.True(slaOutput["bar"])

	meta, err = web.MockGet(app, "/status/sla").JSON(&slaOutput)
	its.Nil(err)
	its.Equal(http.StatusOK, meta.StatusCode)
	its.True(slaOutput["foo"])
	its.True(slaOutput["bar"])

	its.Equal(3, fooCalls)
	its.Equal(3, barCalls)
}

func Test_Controller_getStatusSLA_intermittentFailures(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	var successCalls int32
	var maybeFailureCalls int32
	shouldFail := true
	statusController := NewController(
		OptCheck("success", CheckerFunc(func(_ context.Context) error {
			atomic.AddInt32(&successCalls, 1)
			return nil
		})),
		OptCheck("maybeFailure", CheckerFunc(func(_ context.Context) error {
			atomic.AddInt32(&maybeFailureCalls, 1)
			if shouldFail {
				return fmt.Errorf("oh noes")
			}
			return nil
		})),
	)

	app := web.MustNew()
	app.Register(statusController)

	var slaOutput FreeformResult
	meta, err := web.MockGet(app, "/status/sla").JSON(&slaOutput)
	its.Nil(err)
	its.Equal(http.StatusServiceUnavailable, meta.StatusCode)
	its.True(slaOutput["success"])
	its.False(slaOutput["maybeFailure"])

	shouldFail = false
	meta, err = web.MockGet(app, "/status/sla").JSON(&slaOutput)
	its.Nil(err)
	its.Equal(http.StatusOK, meta.StatusCode)
	its.True(slaOutput["success"])
	its.True(slaOutput["maybeFailure"])

	shouldFail = true
	meta, err = web.MockGet(app, "/status/sla").JSON(&slaOutput)
	its.Nil(err)
	its.Equal(http.StatusServiceUnavailable, meta.StatusCode)
	its.True(slaOutput["success"])
	its.False(slaOutput["maybeFailure"])

	its.Equal(3, successCalls)
	its.Equal(3, maybeFailureCalls)
}

func Test_Controller_getStatusDetails(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	statusController := NewController()

	var shouldFail bool
	interceptor := statusController.Interceptor("test-service")
	action := interceptor.Intercept(ActionerFunc(func(ctx context.Context, args interface{}) (interface{}, error) {
		if shouldFail {
			return nil, fmt.Errorf("oh noes")
		}
		return nil, nil
	}))
	its.Len(statusController.TrackedActions.TrackedActions, 1)

	app := web.MustNew()
	app.Register(statusController)

	var res TrackedActionsResult
	meta, err := web.MockGet(app, "/status/details").JSON(&res)
	its.Nil(err)
	its.Equal(http.StatusOK, meta.StatusCode)
	its.Equal(SignalGreen, res.Status)
	its.Equal(SignalGreen, res.SubSystems["test-service"].Status)

	_, err = action.Action(context.Background(), nil)
	its.Nil(err)

	meta, err = web.MockGet(app, "/status/details").JSON(&res)
	its.Nil(err)
	its.Equal(http.StatusOK, meta.StatusCode)
	its.Equal(SignalGreen, res.Status)

	shouldFail = true

	for x := 0; x < DefaultRedRequestCount; x++ {
		_, err = action.Action(context.Background(), nil)
		its.NotNil(err)
	}

	meta, err = web.MockGet(app, "/status/details").JSON(&res)
	its.Nil(err)
	its.Equal(http.StatusServiceUnavailable, meta.StatusCode)
	its.Equal(SignalRed, res.Status)
	its.Equal(SignalRed, res.SubSystems["test-service"].Status)
}
