/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cron

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func Test_WithJobManager(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	ctx := context.Background()
	ctx = WithJobManager(ctx, New())
	its.NotNil(GetJobManager(ctx))
	its.Nil(GetJobManager(context.Background()))
}

func Test_WithJobScheduler(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	ctx := context.Background()
	ctx = WithJobScheduler(ctx, new(JobScheduler))
	its.NotNil(GetJobScheduler(ctx))
	its.Nil(GetJobScheduler(context.Background()))
}
