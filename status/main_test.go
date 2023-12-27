/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package status

import (
	"context"
	"math"
	"time"
)

func trackedActionConfigDefaults() TrackedActionConfig {
	tac := new(TrackedActionConfig)
	_ = tac.Resolve(context.Background())
	return *tac
}

func trackedActionDefaults() *TrackedAction {
	return &TrackedAction{
		TrackedActionConfig: trackedActionConfigDefaults(),
	}
}

func trackedActionDefaultsWithCounts(requestCount, errorCount int) *TrackedAction {
	ta := &TrackedAction{
		TrackedActionConfig: trackedActionConfigDefaults(),
	}
	for x := 0; x < requestCount; x++ {
		ta.requests = append(ta.requests, RequestInfo{RequestTime: time.Now().UTC()})
	}
	for x := 0; x < errorCount; x++ {
		ta.errors = append(ta.errors, ErrorInfo{Args: "test-endpoint", RequestInfo: RequestInfo{RequestTime: time.Now().UTC()}})
	}
	return ta
}

func lessThanPercentage(target int, percentage float64) (requestCount int, expected float64) {
	requestCount = int(math.Floor(float64(target)/percentage)) - target>>1
	expected = percentage * float64(requestCount)
	return
}

func moreThanPercentage(target int, percentage float64) (requestCount int, expected float64) {
	requestCount = int(math.Ceil(float64(target)/percentage)) + target>>1
	expected = percentage * float64(requestCount)
	return
}
