/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cron

import "context"

type jobManagerKey struct{}

// WithJobManager adds a job manager to a context.
func WithJobManager(ctx context.Context, jm *JobManager) context.Context {
	return context.WithValue(ctx, jobManagerKey{}, jm)
}

// GetJobManager gets a JobManager off a context.
func GetJobManager(ctx context.Context) *JobManager {
	if value := ctx.Value(jobManagerKey{}); value != nil {
		if typed, ok := value.(*JobManager); ok {
			return typed
		}
	}
	return nil
}

type jobSchedulerKey struct{}

// WithJobScheduler adds a job scheduler to a context.
func WithJobScheduler(ctx context.Context, js *JobScheduler) context.Context {
	return context.WithValue(ctx, jobSchedulerKey{}, js)
}

// GetJobScheduler gets a JobScheduler off a context.
func GetJobScheduler(ctx context.Context) *JobScheduler {
	if value := ctx.Value(jobSchedulerKey{}); value != nil {
		if typed, ok := value.(*JobScheduler); ok {
			return typed
		}
	}
	return nil
}
