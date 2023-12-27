/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package status

import (
	"net/http"
	"sort"

	"github.com/blend/go-sdk/async"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

// NewTrackedActionAggregator returns a new tracked action aggregator.
func NewTrackedActionAggregator(trackedActions ...*TrackedAction) *TrackedActionAggregator {
	taa := &TrackedActionAggregator{
		TrackedActions: make(map[string]*TrackedAction),
	}
	for index := range trackedActions {
		trackedAction := trackedActions[index]
		taa.TrackedActions[trackedAction.ServiceName] = trackedAction
	}
	return taa
}

// TrackedActionAggregator aggregates tracker results.
type TrackedActionAggregator struct {
	TrackedActions	map[string]*TrackedAction
	Log		logger.Log
}

// Interceptor returns a new tracked action.
func (taa *TrackedActionAggregator) Interceptor(serviceName string, opts ...TrackedActionOption) async.Interceptor {
	trackedAction := NewTrackedAction(serviceName, opts...)
	taa.TrackedActions[serviceName] = trackedAction
	return trackedAction
}

// Endpoint implements the endpoint.
func (taa *TrackedActionAggregator) Endpoint(servicesToCheck ...string) web.Action {
	return func(r *web.Ctx) web.Result {
		statuses := make(map[string]Info)
		for _, serviceName := range taa.servicesOrDefault(servicesToCheck...) {
			if tracker, ok := taa.TrackedActions[serviceName]; ok {
				statuses[serviceName] = tracker.GetStatus()
			}
		}
		statusCode := http.StatusOK
		status := taa.getSummarySignal(statuses)
		if status != SignalGreen {
			statusCode = http.StatusServiceUnavailable
		}
		return web.JSON.Status(statusCode, TrackedActionsResult{
			Status:		status,
			SubSystems:	statuses,
		})
	}
}

//
// Private / Internal
//

// getSummarySignal implements a this or that (green | red) based on if _any_ of the infos aren't green.
func (taa TrackedActionAggregator) getSummarySignal(statuses map[string]Info) (signal Signal) {
	signal = SignalGreen
	for _, status := range statuses {
		if status.Status != SignalGreen {
			signal = SignalRed
			return
		}
	}
	return
}

// servicesOrDefault returns either the servicesToCheck list if it is set
// or all the keys in the detailed service checks.
func (taa TrackedActionAggregator) servicesOrDefault(servicesToCheck ...string) []string {
	if len(servicesToCheck) > 0 {
		return servicesToCheck
	}
	var output []string
	for key := range taa.TrackedActions {
		output = append(output, key)
	}
	sort.Strings(output)
	return output
}
