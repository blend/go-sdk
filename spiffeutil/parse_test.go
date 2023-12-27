/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package spiffeutil_test

import (
	"testing"

	sdkAssert "github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"

	"github.com/blend/go-sdk/spiffeutil"
)

func TestParse(t *testing.T) {
	assert := sdkAssert.New(t)

	type failureCase struct {
		URI	string
		Message	string
	}
	failures := []failureCase{
		{URI: "https://web.invalid", Message: "Does not match protocol: \"https://web.invalid\""},
		{URI: "spiffe://only.local", Message: "Missing workload identifier: \"spiffe://only.local\""},
		{URI: "spiffe://only.local/", Message: "Missing workload identifier: \"spiffe://only.local/\""},
	}
	for _, fc := range failures {
		pu, err := spiffeutil.Parse(fc.URI)
		assert.Nil(pu)
		assert.True(ex.Is(err, spiffeutil.ErrInvalidURI))
		asEx, ok := err.(*ex.Ex)
		assert.True(ok)
		assert.Equal(fc.Message, asEx.Message)
	}

	// Success.
	pu, err := spiffeutil.Parse("spiffe://cluster.local/ns/blend/sa/quasar")
	expected := &spiffeutil.ParsedURI{TrustDomain: "cluster.local", WorkloadID: "ns/blend/sa/quasar"}
	assert.Equal(expected, pu)
	assert.Nil(err)
}

func TestParseKubernetesWorkloadID(t *testing.T) {
	assert := sdkAssert.New(t)

	type testCase struct {
		WorkloadID	string
		Namespace	string
		ServiceAccount	string
	}
	testCases := []testCase{
		{WorkloadID: "ns/light1/sa/bulb", Namespace: "light1", ServiceAccount: "bulb"},
		{WorkloadID: "xy/light1/sa/bulb"},
		{WorkloadID: "ns/light1/xy/bulb"},
		{WorkloadID: "ns/light1/sa/bulb/extra"},
	}
	for _, tc := range testCases {
		kw, err := spiffeutil.ParseKubernetesWorkloadID(tc.WorkloadID)

		if tc.Namespace == "" {
			assert.True(ex.Is(spiffeutil.ErrNonKubernetesWorkload, err))
			assert.Nil(kw)
		} else {
			assert.Nil(err)
			expected := &spiffeutil.KubernetesWorkload{
				Namespace:	tc.Namespace,
				ServiceAccount:	tc.ServiceAccount,
			}
			assert.Equal(expected, kw)
		}
	}
}
