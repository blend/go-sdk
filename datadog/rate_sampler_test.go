package datadog

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_RateSampler(t *testing.T) {
	assert := assert.New(t)

	sampler := RateSampler(0.25)

	var passed int
	for x := 0; x < 1024; x++ {
		if sampler.Sample(nil) {
			passed++
		}
	}
	assert.True(passed > 240)
	assert.True(passed < 270)
}

func Test_RateSampler_FullOn(t *testing.T) {
	assert := assert.New(t)

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

	sampler := RateSampler(-1)

	var passed int
	for x := 0; x < 1024; x++ {
		if sampler.Sample(nil) {
			passed++
		}
	}
	assert.Zero(passed)
}
