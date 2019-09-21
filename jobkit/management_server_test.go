package jobkit

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/blend/go-sdk/r2"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/web"
)

func TestManagementServerStatic(t *testing.T) {
	assert := assert.New(t)

	_, app := createTestManagementServer()

	meta, err := web.MockGet(app, "/static/js/zepto.min.js").Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
}

func TestManagementServerStatus(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()
	jm.StartAsync()
	defer jm.Stop()

	var status cron.JobManagerStatus
	meta, err := web.MockGet(app, "/status.json").JSON(&status)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal(cron.JobManagerStateRunning, status.State)

	jm.Stop()

	meta, err = web.MockGet(app, "/status.json").JSON(&status)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal(cron.JobManagerStateStopped, status.State)
}

func TestManagementServerIndex(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()
	contents, meta, err := web.MockGet(app, "/").Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	jobName := firstJob(jm).Name()
	assert.Contains(string(contents), fmt.Sprintf("/job/%s", jobName))
	assert.Contains(string(contents), "Show job stats and history")
}

func TestManagementServerSearch(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()
	jobName := firstJob(jm).Name()

	contents, meta, err := web.MockGet(app, "/search", r2.OptQueryValue("selector", "name="+jobName)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Contains(string(contents), fmt.Sprintf("/job/%s", jobName))
	assert.Contains(string(contents), "Show job stats and history")
}

func TestManagementServerSearchInvalidSelector(t *testing.T) {
	assert := assert.New(t)

	_, app := createTestManagementServer()

	meta, err := web.MockGet(app, "/search", r2.OptQueryValue("selector", "~~")).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusBadRequest, meta.StatusCode)
}

func TestManagementServerPause(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()
	jm.StartAsync()
	defer jm.Stop()

	meta, err := web.MockGet(app, "/pause").Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.True(jm.Latch.IsPaused())
}

func TestManagementServerResume(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()
	jm.StartAsync()
	defer jm.Stop()

	jm.Pause()
	assert.True(jm.Latch.IsPaused())

	meta, err := web.MockGet(app, "/resume").Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)

	assert.True(jm.Latch.IsStarted())
}

func TestManagementServerJob(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	job := firstJob(jm)
	assert.NotNil(job)
	jobName := job.Name()
	invocationID := job.History[0].ID

	contents, meta, err := web.MockGet(app, fmt.Sprintf("/job/%s", jobName)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Contains(string(contents), jobName)
	assert.Contains(string(contents), invocationID)

	assert.Contains(string(contents), fmt.Sprintf("/job/%s", jobName))
	assert.NotContains(string(contents), "Show job stats and history")
}

func TestManagementServerJobNotFound(t *testing.T) {
	assert := assert.New(t)

	_, app := createTestManagementServer()

	meta, err := web.MockGet(app, fmt.Sprintf("/job/%s", uuid.V4().String())).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusNotFound, meta.StatusCode)
}

func TestManagementServerJobDisable(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	job, err := jm.Job("test1")
	assert.Nil(err)
	assert.NotNil(job)
	jobName := job.Name()
	assert.False(job.Disabled)

	meta, err := web.MockGet(app, fmt.Sprintf("/job.disable/%s", jobName)).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)

	assert.True(job.Disabled)
}

func TestManagementServerJobEnable(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	job, err := jm.Job("test1")
	assert.Nil(err)
	assert.NotNil(job)
	jobName := job.Name()
	job.Disabled = true

	meta, err := web.MockGet(app, fmt.Sprintf("/job.enable/%s", jobName)).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.False(job.Disabled)
}

func TestManagementServerJobRun(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	job, err := jm.Job("test1")
	assert.Nil(err)
	assert.NotNil(job)
	jobName := job.Name()

	meta, err := web.MockGet(app, fmt.Sprintf("/job.run/%s", jobName)).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.NotNil(job.Last)
}

func TestManagementServerJobCancel(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	called := make(chan struct{})
	cancelled := make(chan struct{})

	job := cron.NewJob(cron.OptJobName("cancel-test"), cron.OptJobAction(func(ctx context.Context) error {
		close(called)
		<-ctx.Done()
		close(cancelled)
		return nil
	}))
	jm.LoadJobs(job)

	meta, err := web.MockGet(app, fmt.Sprintf("/job.run/%s", job.Name())).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)

	<-called

	meta, err = web.MockGet(app, fmt.Sprintf("/job.cancel/%s", job.Name())).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)

	<-cancelled
}

func TestManagementServerJobInvocation(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	job := firstJob(jm)
	assert.NotNil(job)

	jobName := job.Name()
	invocationID := job.History[0].ID

	contents, meta, err := web.MockGet(app, fmt.Sprintf("/job.invocation/%s/%s", jobName, invocationID)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode, string(contents))
	assert.Contains(string(contents), jobName)
	assert.Contains(string(contents), invocationID)
}

func TestManagementServerJobInvocationCurrent(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	job := firstJob(jm)
	assert.NotNil(job)

	jobName := job.Name()
	invocationID := job.Current.ID

	contents, meta, err := web.MockGet(app, fmt.Sprintf("/job.invocation/%s/current", jobName)).Bytes()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode, string(contents))
	assert.Contains(string(contents), jobName)
	assert.Contains(string(contents), invocationID)
}

func TestManagementServerJobInvocationNotFound(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	meta, err := web.MockGet(app, fmt.Sprintf("/job.invocation/%s/%s", uuid.V4().String(), uuid.V4().String())).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusNotFound, meta.StatusCode)

	job := firstJob(jm)
	assert.NotNil(job)
	jobName := job.Name()

	meta, err = web.MockGet(app, fmt.Sprintf("/job.invocation/%s/%s", jobName, uuid.V4().String())).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusNotFound, meta.StatusCode)
}

//
// api tests
//

func TestManagementServerAPIJobs(t *testing.T) {
	assert := assert.New(t)

	_, app := createTestManagementServer()
	var jobs []cron.JobSchedulerStatus
	meta, err := web.MockGet(app, "/api/jobs").JSON(&jobs)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.NotEmpty(jobs)
}

func TestManagementServerAPIJobsRunning(t *testing.T) {
	assert := assert.New(t)

	_, app := createTestManagementServer()
	var jobs map[string]cron.JobInvocation
	meta, err := web.MockGet(app, "/api/jobs.running").JSON(&jobs)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.NotEmpty(jobs)
}

func TestManagementServerAPIJobsStats(t *testing.T) {
	assert := assert.New(t)

	_, app := createTestManagementServer()
	var jobs []cron.JobStats
	meta, err := web.MockGet(app, "/api/jobs.stats").JSON(&jobs)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.NotEmpty(jobs)
}

func TestManagementServerAPIJob(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	job := firstJob(jm)
	assert.NotNil(job)
	jobName := job.Name()

	var js cron.JobSchedulerStatus
	meta, err := web.MockGet(app, fmt.Sprintf("/api/job/%s", jobName)).JSON(&js)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal(jobName, js.Name)
}

func TestManagementServerAPIJobNotFound(t *testing.T) {
	assert := assert.New(t)

	_, app := createTestManagementServer()
	meta, err := web.MockGet(app, fmt.Sprintf("/api/job/%s", uuid.V4().String())).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusNotFound, meta.StatusCode)
}

func TestManagementServerAPIJobStats(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	job := firstJob(jm)
	assert.NotNil(job)
	jobName := job.Name()

	var js cron.JobStats
	meta, err := web.MockGet(app, fmt.Sprintf("/api/job.stats/%s", jobName)).JSON(&js)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.Equal(jobName, js.JobName)
}

func TestManagementServerAPIPause(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()
	jm.StartAsync()
	defer jm.Stop()

	meta, err := web.MockPost(app, "/api/pause", nil).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.True(jm.Latch.IsPaused())
}

func TestManagementServerAPIResume(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()
	jm.StartAsync()
	defer jm.Stop()

	jm.Pause()
	assert.True(jm.Latch.IsPaused())

	meta, err := web.MockPost(app, "/api/resume", nil).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.True(jm.Latch.IsStarted())
}

func TestManagementServerAPIJobRun(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	job, err := jm.Job("test1")
	assert.Nil(err)
	assert.NotNil(job)
	jobName := job.Name()

	var ji cron.JobInvocation
	meta, err := web.MockPost(app, fmt.Sprintf("/api/job.run/%s", jobName), nil).JSON(&ji)
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.NotEmpty(ji.ID)
	assert.False(ji.Started.IsZero())
	assert.Equal("test1", ji.JobName)
}

func TestManagementServerAPIJobCancel(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	called := make(chan struct{})
	cancelled := make(chan struct{})

	job := cron.NewJob(cron.OptJobName("cancel-test"), cron.OptJobAction(func(ctx context.Context) error {
		close(called)
		<-ctx.Done()
		close(cancelled)
		return nil
	}))
	jm.LoadJobs(job)

	meta, err := web.MockPost(app, fmt.Sprintf("/api/job.run/%s", job.Name()), nil).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)

	<-called

	meta, err = web.MockPost(app, fmt.Sprintf("/api/job.cancel/%s", job.Name()), nil).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)

	<-cancelled
}

func TestManagementServerAPIJobDisable(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	job, err := jm.Job("test1")
	assert.Nil(err)
	assert.NotNil(job)
	jobName := job.Name()
	assert.False(job.Disabled)

	meta, err := web.MockPost(app, fmt.Sprintf("/api/job.disable/%s", jobName), nil).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)

	assert.True(job.Disabled)
}

func TestManagementServerAPIJobEnable(t *testing.T) {
	assert := assert.New(t)

	jm, app := createTestManagementServer()

	job, err := jm.Job("test1")
	assert.Nil(err)
	assert.NotNil(job)
	jobName := job.Name()
	job.Disabled = true

	meta, err := web.MockPost(app, fmt.Sprintf("/api/job.enable/%s", jobName), nil).Discard()
	assert.Nil(err)
	assert.Equal(http.StatusOK, meta.StatusCode)
	assert.False(job.Disabled)
}
