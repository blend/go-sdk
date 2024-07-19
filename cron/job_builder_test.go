/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
