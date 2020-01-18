package cron

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/uuid"
)

func TestJobInvocationJSON(t *testing.T) {
	assert := assert.New(t)

	test := JobInvocation{
		ID:        uuid.V4().String(),
		JobName:   uuid.V4().String(),
		Started:   time.Date(2019, 9, 4, 12, 11, 10, 9, time.UTC),
		Complete:  time.Date(2019, 9, 6, 12, 11, 10, 9, time.UTC),
		Cancelled: time.Date(2019, 9, 7, 12, 11, 10, 9, time.UTC),
		Timeout:   time.Date(2019, 9, 8, 12, 11, 10, 9, time.UTC),
		Status:    JobInvocationStatusComplete,
		Err:       ex.New("this is a test"),
	}

	contents, err := json.Marshal(test)
	assert.Nil(err)
	assert.NotEmpty(contents)

	var verify JobInvocation
	assert.Nil(json.Unmarshal(contents, &verify))
	assert.Equal(test.ID, verify.ID)
	assert.Equal(test.JobName, verify.JobName)
	assert.Equal(test.Started, verify.Started)
	assert.Equal(test.Complete, verify.Complete)
	assert.Equal(test.Cancelled, verify.Cancelled)
	assert.Equal(test.Timeout, verify.Timeout)
	assert.Equal(test.Status, verify.Status)

	assert.NotNil(verify.Err)
	assert.Contains(verify.Err.Error(), "this is a test")
}
