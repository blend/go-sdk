/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package status

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_TrackedActionAggregator_getSummarySignal(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	empty := map[string]Info{}
	its.Equal(SignalGreen, new(TrackedActionAggregator).getSummarySignal(empty))

	ok := map[string]Info{
		"foo":	{Status: SignalGreen},
		"bar":	{Status: SignalGreen},
		"baz":	{Status: SignalGreen},
	}
	its.Equal(SignalGreen, new(TrackedActionAggregator).getSummarySignal(ok))

	mixed := map[string]Info{
		"foo":	{Status: SignalGreen},
		"bar":	{Status: SignalYellow},
		"baz":	{Status: SignalRed},
	}
	its.Equal(SignalRed, new(TrackedActionAggregator).getSummarySignal(mixed))

	bad := map[string]Info{
		"foo":	{Status: SignalRed},
		"bar":	{Status: SignalRed},
		"baz":	{Status: SignalRed},
	}
	its.Equal(SignalRed, new(TrackedActionAggregator).getSummarySignal(bad))
}

func Test_TrackedActionAggregator_servicesOrDefault(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	taa := TrackedActionAggregator{
		TrackedActions: map[string]*TrackedAction{
			"foo":	nil,
			"bar":	nil,
			"baz":	nil,
		},
	}

	defaults := []string{"bar", "baz", "foo"}	// map keys get sorted
	its.Equal(defaults, taa.servicesOrDefault())

	servicesToCheck := []string{"alpha", "bravo", "charlie"}
	its.Equal(servicesToCheck, taa.servicesOrDefault(servicesToCheck...))
}
