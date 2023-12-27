/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package datadog

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_RateSampler(t *testing.T) {
	assert := assert.New(t)

	// 25% sample (25% pass at infinity)
	sampler := RateSampler(0.25)

	var passed int
	for x := 0; x < 102400; x++ {
		if sampler.Sample(nil) {
			passed++
		}
	}

	// Since we are sampling incrementally, we cannot guarantee that we would
	// sample exactly 25% of the population, so we add a buffer here to account
	// for the errors. The ratio of errors gets smaller the larger the population.
	expected := int(102400 * 0.25)
	buffer := 1024
	assert.True(passed > expected-buffer, fmt.Sprint(passed))
	assert.True(passed < expected+buffer, fmt.Sprint(passed))
}

func Test_RateSampler_FullOn(t *testing.T) {
	assert := assert.New(t)

	// 100% sample (all pass)
	sampler := RateSampler(1)

	var passed int
	for x := 0; x < 1024; x++ {
		if sampler.Sample(nil) {
			passed++
		}
	}
	assert.Equal(passed, 1024)
}

func Test_RateSampler_FullOff(t *testing.T) {
	assert := assert.New(t)

	// 0% sample (none passes)
	sampler := RateSampler(0)

	var passed int
	for x := 0; x < 1024; x++ {
		if sampler.Sample(nil) {
			passed++
		}
	}
	assert.Zero(passed)
}

func Test_RateSampler_FullOff_Negative(t *testing.T) {
	assert := assert.New(t)

	// Negative sampling rate (none passes)
	sampler := RateSampler(-1)

	var passed int
	for x := 0; x < 1024; x++ {
		if sampler.Sample(nil) {
			passed++
		}
	}
	assert.Zero(passed)
}
