/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package cron

import (
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
)

func TestJobInvocationElapsed(t *testing.T) {
	assert := assert.New(t)

	started := time.Now().UTC()

	assert.Equal(200*time.Millisecond, (&JobInvocation{
		Started:	started,
		Complete:	started.Add(200 * time.Millisecond),
	}).Elapsed())

	assert.NotZero((&JobInvocation{
		Started: started,
	}).Elapsed())
}

func TestJobInvocationClone(t *testing.T) {
	assert := assert.New(t)

	ts := time.Now().UTC()
	ji := &JobInvocation{
		ID:		NewJobInvocationID(),
		JobName:	uuid.V4().String(),
		Started:	ts,
		Complete:	ts.Add(100 * time.Millisecond),
		Err:		fmt.Errorf("this is a test"),
		Status:		JobInvocationStatusErrored,
		Parameters:	map[string]string{"foo": "bar", "example-string": "dog"},
		State:		"this is also a test",
		Cancel:		func() {},
	}
	cloned := ji.Clone()
	assert.Equal(ji.ID, cloned.ID)
	assert.Equal(ji.JobName, cloned.JobName)
	assert.Equal(ji.Started, cloned.Started)
	assert.Equal(ji.Complete, cloned.Complete)
	assert.Equal(ji.Err, cloned.Err)
	assert.Equal(ji.Status, cloned.Status)
	assert.Equal(ji.Parameters, cloned.Parameters)
	assert.Equal(ji.State, cloned.State)
	assert.NotNil(cloned.Cancel)
}
