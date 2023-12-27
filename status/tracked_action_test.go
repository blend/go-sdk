/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package status

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func Test_TrackedAction_Wrap(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	var shouldError, shouldPanic bool
	action := ActionerFunc(func(_ context.Context, _ interface{}) (interface{}, error) {
		if shouldPanic {
			panic(fmt.Errorf("this is a panic"))
		}
		if shouldError {
			return nil, fmt.Errorf("this is an error")
		}
		return "ok!", nil
	})
	ta := NewTrackedAction("test-service")
	tracked := ta.Intercept(action)

	// should yield ok!
	res, err := tracked.Action(context.Background(), "test-resource")
	its.Nil(err)
	its.Equal("ok!", res)
	its.Len(ta.errors, 0)
	its.Len(ta.requests, 1)

	shouldError = true
	res, err = tracked.Action(context.Background(), "test-resource")
	its.Equal("this is an error", ex.ErrClass(err).Error())
	its.Nil(res)
	its.Len(ta.errors, 1)
	its.Len(ta.requests, 1)

	shouldPanic = true
	shouldError = false
	res, err = tracked.Action(context.Background(), "test-resource")
	its.Equal("this is a panic", ex.ErrClass(err).Error())
	its.Nil(res)
	its.Len(ta.errors, 2)
	its.Len(ta.requests, 1)

	// push "now" forward by the expiration, the next wrapped call should clear the history
	shouldPanic = false
	shouldError = false
	ta.nowProvider = func() time.Time { return time.Now().UTC().Add(ta.ExpirationOrDefault()) }
	res, err = tracked.Action(context.Background(), "test-resource")
	its.Nil(err)
	its.Equal("ok!", res)
	its.Len(ta.errors, 0)
	its.Len(ta.requests, 1)
}

func Test_TrackedAction_GetStatus(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	now := time.Now().UTC()
	ta := &TrackedAction{
		ServiceName:	"test-tracked-action",
		nowProvider:	func() time.Time { return now },
		TrackedActionConfig: TrackedActionConfig{
			Expiration: 5 * time.Second,
		},
		requests: []RequestInfo{
			{now},
			{now.Add(-1 * time.Second)},
			{now.Add(-2 * time.Second)},
			{now.Add(-3 * time.Second)},
			{now.Add(-4 * time.Second)},
			{now.Add(-5 * time.Second)},
			{now.Add(-6 * time.Second)},
			{now.Add(-7 * time.Second)},
		},
		errors: []ErrorInfo{
			{Args: "test-resource-0", RequestInfo: RequestInfo{now.Add(-2 * time.Second)}},
			{Args: "test-resource-1", RequestInfo: RequestInfo{now.Add(-3 * time.Second)}},
			{Args: "test-resource-0", RequestInfo: RequestInfo{now.Add(-4 * time.Second)}},
			{Args: "test-resource-1", RequestInfo: RequestInfo{now.Add(-5 * time.Second)}},
			{Args: "test-resource-0", RequestInfo: RequestInfo{now.Add(-6 * time.Second)}},
			{Args: "test-resource-1", RequestInfo: RequestInfo{now.Add(-7 * time.Second)}},
			{Args: "test-resource-1", RequestInfo: RequestInfo{now.Add(-8 * time.Second)}},
		},
	}

	status := ta.GetStatus()
	its.Equal("test-tracked-action", status.Name)
	its.Equal(SignalRed, status.Status)
	its.Equal(5, status.Details.RequestCount)
	its.Equal(3, status.Details.ErrorCount)
	its.Len(status.Details.ErrorBreakdown, 2)
	its.Equal(2, status.Details.ErrorBreakdown["test-resource-0"])
	its.Equal(1, status.Details.ErrorBreakdown["test-resource-1"])
}

func Test_TrackedAction_getStatusSignalUnsafe(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	testCases := [...]struct {
		RequestCount	int
		ErrorCount	int
		Expected	Signal
		Message		string
	}{
		{0, 0, SignalGreen, "should return green with no requests"},
		{0, 10, SignalYellow, "should return yellow when 10 requests fail"},
		{0, 50, SignalRed, "should return red when 50 requests fail"},
		{2200, 10, SignalGreen, "should use percentages when count is high enough"},
	}
	for _, tc := range testCases {
		its.Equal(tc.Expected, trackedActionDefaultsWithCounts(tc.RequestCount, tc.ErrorCount).getStatusSignalUnsafe(), tc.Message)
	}
}

func Test_TrackedAction_cleanOldRequestsUnsafe(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	now := time.Now().UTC()
	ta := &TrackedAction{
		nowProvider:	func() time.Time { return now },
		TrackedActionConfig: TrackedActionConfig{
			Expiration: 5 * time.Second,
		},
		requests: []RequestInfo{
			{now},
			{now.Add(-1 * time.Second)},
			{now.Add(-2 * time.Second)},
			{now.Add(-3 * time.Second)},
			{now.Add(-4 * time.Second)},
			{now.Add(-5 * time.Second)},
			{now.Add(-6 * time.Second)},
			{now.Add(-7 * time.Second)},
		},
		errors: []ErrorInfo{
			{RequestInfo: RequestInfo{now.Add(-3 * time.Second)}},
			{RequestInfo: RequestInfo{now.Add(-4 * time.Second)}},
			{RequestInfo: RequestInfo{now.Add(-5 * time.Second)}},
			{RequestInfo: RequestInfo{now.Add(-6 * time.Second)}},
			{RequestInfo: RequestInfo{now.Add(-7 * time.Second)}},
			{RequestInfo: RequestInfo{now.Add(-8 * time.Second)}},
		},
	}

	ta.cleanOldRequestsUnsafe()
	its.Len(ta.requests, 5)
	its.Len(ta.errors, 2)
}

func Test_TrackedAction_redErrorCount(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	lessThanRequestCount, _ := lessThanPercentage(DefaultRedRequestCount, DefaultRedRequestPercentage)
	moreThanRequestCount, moreThanExpected := moreThanPercentage(DefaultRedRequestCount, DefaultRedRequestPercentage)
	testCases := [...]struct {
		RequestCount	int
		Expected	float64
	}{
		{lessThanRequestCount, DefaultRedRequestCount},
		{moreThanRequestCount, moreThanExpected},
	}

	for _, tc := range testCases {
		its.Equal(tc.Expected, trackedActionDefaults().redErrorCount(tc.RequestCount), fmt.Sprintf("requestCount: %d", tc.RequestCount))
	}
}

func Test_TrackedAction_yellowErrorCount(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	lessThanRequestCount, _ := lessThanPercentage(DefaultYellowRequestCount, DefaultYellowRequestPercentage)
	moreThanRequestCount, moreThanExpected := moreThanPercentage(DefaultYellowRequestCount, DefaultYellowRequestPercentage)
	testCases := [...]struct {
		RequestCount	int
		Expected	float64
	}{
		{lessThanRequestCount, DefaultYellowRequestCount},
		{moreThanRequestCount, moreThanExpected},
	}

	for _, tc := range testCases {
		its.Equal(tc.Expected, trackedActionDefaults().yellowErrorCount(tc.RequestCount), fmt.Sprintf("requestCount: %d", tc.RequestCount))
	}
}

func Test_TrackedAction_now(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	now := time.Date(2021, 06, 13, 11, 50, 0, 0, time.UTC)
	withProvider := TrackedAction{
		nowProvider: func() time.Time {
			return now
		},
	}

	its.Equal(now, withProvider.now())
	its.NotEqual(now, new(TrackedAction).now())
}
