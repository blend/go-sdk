/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cron

import (
	"context"
	"time"

	"github.com/blend/go-sdk/uuid"
)

// NewJobInvocation returns a new job invocation.
func NewJobInvocation(jobName string) *JobInvocation {
	return &JobInvocation{
		ID:      NewJobInvocationID(),
		Status:  JobInvocationStatusIdle,
		JobName: jobName,
	}
}

type contextKeyJobInvocation struct{}

// WithJobInvocation adds job invocation to a context.
func WithJobInvocation(ctx context.Context, ji *JobInvocation) context.Context {
	return context.WithValue(ctx, contextKeyJobInvocation{}, ji)
}

// GetJobInvocation gets the job invocation from a given context.
func GetJobInvocation(ctx context.Context) *JobInvocation {
	if value := ctx.Value(contextKeyJobInvocation{}); value != nil {
		if typed, ok := value.(*JobInvocation); ok {
			return typed
		}
	}
	return nil
}

// NewJobInvocationID returns a new pseudo-unique job invocation identifier.
func NewJobInvocationID() string {
	return uuid.V4().String()
}

// JobInvocation is metadata for a job invocation (or instance of a job running).
type JobInvocation struct {
	ID      string `json:"id"`
	JobName string `json:"jobName"`

	Started  time.Time `json:"started"`
	Complete time.Time `json:"complete"`
	Err      error     `json:"err"`

	Parameters JobParameters       `json:"parameters"`
	Status     JobInvocationStatus `json:"status"`
	State      interface{}         `json:"-"`

	Cancel context.CancelFunc `json:"-"`
}

// Elapsed returns the elapsed time for the invocation.
func (ji *JobInvocation) Elapsed() time.Duration {
	if !ji.Complete.IsZero() {
		return ji.Complete.Sub(ji.Started)
	}
	if !ji.Started.IsZero() {
		return Now().Sub(ji.Started)
	}
	return 0
}

// Clone clones the job invocation.
func (ji *JobInvocation) Clone() *JobInvocation {
	return &JobInvocation{
		ID:      ji.ID,
		JobName: ji.JobName,

		Started:  ji.Started,
		Complete: ji.Complete,
		Err:      ji.Err,

		Parameters: ji.Parameters,
		Status:     ji.Status,
		State:      ji.State,

		Cancel: ji.Cancel,
	}
}
