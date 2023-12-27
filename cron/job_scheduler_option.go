/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package cron

import (
	"context"

	"github.com/blend/go-sdk/logger"
)

// JobSchedulerOption is an option for job schedulers.
type JobSchedulerOption func(*JobScheduler)

// OptJobSchedulerTracer sets the job scheduler tracer.
func OptJobSchedulerTracer(tracer Tracer) JobSchedulerOption {
	return func(js *JobScheduler) { js.Tracer = tracer }
}

// OptJobSchedulerLog sets the job scheduler logger.
func OptJobSchedulerLog(log logger.Log) JobSchedulerOption {
	return func(js *JobScheduler) { js.Log = log }
}

// OptJobSchedulerBaseContext sets the job scheduler BaseContext.
func OptJobSchedulerBaseContext(ctx context.Context) JobSchedulerOption {
	return func(js *JobScheduler) { js.BaseContext = ctx }
}
