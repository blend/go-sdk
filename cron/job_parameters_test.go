/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

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
