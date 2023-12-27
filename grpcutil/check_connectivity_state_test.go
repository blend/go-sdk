/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package grpcutil

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/connectivity"

	"github.com/blend/go-sdk/assert"
)

type mockGetConnectionState connectivity.State

func (m mockGetConnectionState) GetConnectionState() connectivity.State {
	return connectivity.State(m)
}

type mockGetConnectionStateMany chan connectivity.State

func (m mockGetConnectionStateMany) GetConnectionState() connectivity.State {
	return connectivity.State(<-m)
}

func Test_CheckConnectivityState(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	checker := CheckConnectivityState(
		mockGetConnectionState(connectivity.Ready),
	)
	err := checker.Check(context.Background())
	its.Nil(err)
}

func Test_CheckConnectivityState_failure(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	checker := CheckConnectivityState(
		mockGetConnectionState(connectivity.TransientFailure),
		OptRetryCheckConnectivityStateRetryBackoff(time.Microsecond),
		OptRetryCheckConnectivityStateRetryTimeout(time.Millisecond),
	)
	err := checker.Check(context.Background())
	its.NotNil(err)
}

func Test_CheckConnectivityState_retry_success(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	states := mockGetConnectionStateMany(make(chan connectivity.State, 5))
	states <- connectivity.Connecting
	states <- connectivity.Connecting
	states <- connectivity.Connecting
	states <- connectivity.Connecting
	states <- connectivity.Ready

	checker := CheckConnectivityState(
		mockGetConnectionState(connectivity.Ready),
		OptRetryCheckConnectivityStateRetryBackoff(time.Microsecond),
		OptRetryCheckConnectivityStateRetryTimeout(time.Millisecond),
	)
	err := checker.Check(context.Background())
	its.Nil(err)
}
