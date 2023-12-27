/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package cron

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_JobParametersContext(t *testing.T) {
	assert := assert.New(t)

	final := GetJobParameterValues(WithJobParameterValues(context.Background(), JobParameters{
		"foo":	"bar",
		"buzz":	"fuzz",
	}))
	assert.Equal("bar", final["foo"])
	assert.Equal("fuzz", final["buzz"])

	assert.Empty(GetJobParameterValues(context.Background()))
}
