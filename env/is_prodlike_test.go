/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package env_test

import (
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/env"
	"github.com/blend/go-sdk/uuid"
)

func TestIsProdlike(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		Input		string
		Expected	bool
	}{
		{Input: env.ServiceEnvDev, Expected: false},
		{Input: env.ServiceEnvCI, Expected: false},
		{Input: env.ServiceEnvTest, Expected: false},
		{Input: env.ServiceEnvSandbox, Expected: false},
		{Input: env.ServiceEnvPreprod, Expected: true},
		{Input: env.ServiceEnvBeta, Expected: true},
		{Input: env.ServiceEnvProd, Expected: true},
		{Input: uuid.V4().String(), Expected: true},
		{Expected: true},
	}

	for _, testCase := range testCases {
		assert.Equal(testCase.Expected, env.IsProdlike(testCase.Input), fmt.Sprintf("failed for: %s", testCase.Input))
	}
}

func TestIsProduction(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		Input		string
		Expected	bool
	}{
		{Input: env.ServiceEnvDev, Expected: false},
		{Input: env.ServiceEnvCI, Expected: false},
		{Input: env.ServiceEnvTest, Expected: false},
		{Input: env.ServiceEnvSandbox, Expected: false},
		{Input: env.ServiceEnvPreprod, Expected: true},
		{Input: env.ServiceEnvBeta, Expected: false},
		{Input: env.ServiceEnvProd, Expected: true},
		{Input: uuid.V4().String(), Expected: false},
		{Expected: false},
	}

	for _, testCase := range testCases {
		assert.Equal(testCase.Expected, env.IsProduction(testCase.Input), fmt.Sprintf("failed for: %s", testCase.Input))
	}
}
