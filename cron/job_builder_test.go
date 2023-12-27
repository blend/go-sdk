/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package cron

import (
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestJobBuilder(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(NewJob())
	assert.Equal("test_job", NewJob(OptJobName("test_job")).Name())
}
