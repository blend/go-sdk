/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package status

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/ex"
)

var (
	_ async.Interceptor = (*TrackedAction)(nil)
)

// NewTrackedAction returns a new tracked action.
func NewTrackedAction(serviceName string, opts ...TrackedActionOption) *TrackedAction {
	ta := &TrackedAction{
		ServiceName: serviceName,
	}
	_ = (&ta.TrackedActionConfig).Resolve(context.Background())
	for _, opt := range opts {
		opt(ta)
	}
	return ta
}

// TrackedActionOption mutates a tracked action.
type TrackedActionOption func(*TrackedAction)

// OptTrackedActionConfig sets the tracked action config.
func OptTrackedActionConfig(cfg TrackedActionConfig) TrackedActionOption {
	return func(ta *TrackedAction) {
		ta.TrackedActionConfig = cfg
	}
}

// TrackedAction is a wrapper for action that tracks a rolling
// window of history based on the configured expiration.
type TrackedAction struct {
	TrackedActionConfig
	sync.Mutex

	ServiceName	string

	nowProvider	func() time.Time
	errors		[]ErrorInfo
	requests	[]RequestInfo
}

// Intercept implements async.Interceptor.
func (t *TrackedAction) Intercept(action Actioner) Actioner {
	return ActionerFunc(func(ctx context.Context, args interface{}) (output interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = ex.Append(err, ex.New(r))
			}
			t.CleanOldRequests()
			if err != nil {
				t.AddErroredRequest(args)
			} else {
				t.AddSuccessfulRequest()
			}
		}()
		output, err = action.Action(ctx, args)
		return
	})
}

// GetStatus gets the status for the tracker.
//
// It is safe to call concurrently from multiple goroutines.
func (t *TrackedAction) GetStatus() (info Info) {
	t.Lock()
	defer t.Unlock()

	t.cleanOldRequestsUnsafe()
	info.Name = t.ServiceName
	info.Status = t.getStatusSignalUnsafe()

	errorBreakdown := make(map[string]int)
	if info.Status == SignalYellow || info.Status == SignalRed {
		for _, errorInfo := range t.errors {
			errorBreakdown[t.formatArgs(errorInfo.Args)]++
		}
	}
	info.Details = Details{
		ErrorCount:	len(t.errors),
		RequestCount:	len(t.requests),
		ErrorBreakdown:	errorBreakdown,
	}
	return
}

// GetStatusSignal returns the current status signal.
//
// It is safe to call concurrently from multiple goroutines.
func (t *TrackedAction) GetStatusSignal() (status Signal) {
	t.Lock()
	status = t.getStatusSignalUnsafe()
	t.Unlock()
	return
}

// AddErroredRequest adds an errored request.
//
// It is safe to call concurrently from multiple goroutines.
func (t *TrackedAction) AddErroredRequest(args interface{}) {
	t.Lock()
	defer t.Unlock()
	t.errors = append(t.errors, ErrorInfo{
		Args:	args,
		RequestInfo: RequestInfo{
			RequestTime: t.now(),
		},
	})
}

// AddSuccessfulRequest adds a successful request.
//
// It is safe to call concurrently from multiple goroutines.
func (t *TrackedAction) AddSuccessfulRequest() {
	t.Lock()
	defer t.Unlock()
	t.requests = append(t.requests, RequestInfo{RequestTime: t.now()})
}

// CleanOldRequests is an action delegate that removes expired requests
// from the tracker
//
// It is safe to call concurrently from multiple goroutines.
func (t *TrackedAction) CleanOldRequests() {
	t.Lock()
	defer t.Unlock()
	t.cleanOldRequestsUnsafe()
}

//
// Private - Internal
//

func (t *TrackedAction) formatArgs(args interface{}) string {
	switch typed := args.(type) {
	case string:
		return typed
	case []byte:
		return string(typed)
	case []rune:
		return string(typed)
	case fmt.Stringer:
		return typed.String()
	default:
		return "unknown"
	}
}

// getStatusSignalUnsafe gets the specific signal (green, yellow, or red)
// for the tracker.
func (t *TrackedAction) getStatusSignalUnsafe() (status Signal) {
	status = SignalGreen
	requestCount := len(t.requests)
	errorCount := float64(len(t.errors))
	if errorCount >= t.redErrorCount(requestCount) {
		status = SignalRed
	} else if errorCount >= t.yellowErrorCount(requestCount) {
		status = SignalYellow
	}
	return status
}

func (t *TrackedAction) cleanOldRequestsUnsafe() {
	nowUTC := t.now()
	var filteredErrors []ErrorInfo
	for _, errorInfo := range t.errors {
		if nowUTC.Sub(errorInfo.RequestTime) < t.ExpirationOrDefault() {
			filteredErrors = append(filteredErrors, errorInfo)
		}
	}

	t.errors = filteredErrors
	var filteredRequests []RequestInfo
	for _, requestInfo := range t.requests {
		if nowUTC.Sub(requestInfo.RequestTime) < t.ExpirationOrDefault() {
			filteredRequests = append(filteredRequests, requestInfo)
		}
	}
	t.requests = filteredRequests
}

// redErrorCount returns the expected threshold for what is
// considered a "red" signal status based on either the baseline `RedRequestCount`
// or the RedRequestPercentage applied to the current request count.
//
// It is meant to scale the threshold to the volume of the calls
// to the tracked action.
func (t *TrackedAction) redErrorCount(requestCount int) float64 {
	return math.Max(
		float64(t.RedRequestCount),
		t.RedRequestPercentage*float64(requestCount),
	)
}

// yellowErrorCount returns the expected threshold for what is
// considered a "yellow" signal status based on either the baseline `YellowRequestCount`
// or the YellowRequestPercentage applied to the current request count.
//
// It is meant to scale the threshold to the volume of the calls
// to the tracked action
func (t *TrackedAction) yellowErrorCount(requestCount int) float64 {
	return math.Max(
		float64(t.YellowRequestCount),
		t.YellowRequestPercentage*float64(requestCount),
	)
}

// now returns the current time.
func (t *TrackedAction) now() time.Time {
	if t.nowProvider != nil {
		return t.nowProvider()
	}
	return time.Now().UTC()
}
