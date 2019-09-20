package jobkit

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/selector"
	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/web"
	"github.com/blend/go-sdk/webutil"

	"github.com/blend/go-sdk/jobkit/static"
	"github.com/blend/go-sdk/jobkit/views"
)

// NewServer returns a new management server that lets you
// trigger jobs or look at job statuses via. a json api.
func NewServer(jm *cron.JobManager, cfg Config, options ...web.Option) *web.App {
	options = append([]web.Option{web.OptConfig(cfg.Web)}, options...)
	app := web.MustNew(options...)
	app.Register(ManagementServer{Cron: jm, Config: cfg})
	return app
}

// ManagementServer is the jobkit management server.
type ManagementServer struct {
	Config Config
	Cron   *cron.JobManager
}

// Register registers the management server.
func (ms ManagementServer) Register(app *web.App) {
	if ms.Config.UseViewFilesOrDefault() {
		app.Views.LiveReload = true
		app.Views.AddPaths(ms.ViewPaths()...)
	} else {
		app.Views.LiveReload = false
		for _, viewPath := range ms.ViewPaths() {
			vf, err := views.GetBinaryAsset(viewPath)
			if err != nil {
				panic(err)
			}
			contents, err := vf.Contents()
			if err != nil {
				panic(err)
			}
			app.Views.AddLiterals(string(contents))
		}
	}
	app.DefaultMiddleware = append(app.DefaultMiddleware, ms.addContextStateConfig)
	app.PanicAction = func(r *web.Ctx, err interface{}) web.Result {
		return r.Views.InternalError(ex.New(err))
	}

	// web specific routes
	app.GET("/status.json", ms.getStatus)
	app.GET("/static/*filepath", ms.getStatic)

	// manager routes
	app.GET("/", ms.getIndex)
	app.GET("/search", ms.getSearch)
	app.GET("/pause", ms.getPause)
	app.GET("/resume", ms.getResume)

	// job routes
	app.GET("/job/:jobName", ms.getJob)
	app.GET("/job.run/:jobName", ms.getJobRun)
	app.GET("/job.enable/:jobName", ms.getJobEnable)
	app.GET("/job.disable/:jobName", ms.getJobDisable)
	app.GET("/job.cancel/:jobName", ms.getJobCancel)

	// invocation routes
	app.GET("/job.invocation/:jobName/:id", ms.getJobInvocation)

	// api routes
	app.POST("/api/pause", ms.postAPIPause)
	app.POST("/api/resume", ms.postAPIResume)
	app.GET("/api/jobs", ms.getAPIJobs)
	app.GET("/api/jobs.running", ms.getAPIJobsRunning)
	app.GET("/api/jobs.stats", ms.getAPIJobStats)
	app.GET("/api/job/:jobName", ms.getAPIJob)
	app.GET("/api/job.stats/:jobName", ms.getAPIJobStats)
	app.POST("/api/job.run/:jobName", ms.postAPIJobRun)
	app.POST("/api/job.cancel/:jobName", ms.postAPIJobCancel)
	app.POST("/api/job.disable/:jobName", ms.postAPIJobDisable)
	app.POST("/api/job.enable/:jobName", ms.postAPIJobEnable)
	app.GET("/api/job.invocation/:jobName/:id", ms.getAPIJobInvocation)
	app.GET("/api/job.invocation.output/:jobName/:id", ms.getAPIJobInvocationOutput)
	app.GET("/api/job.invocation.output.stream/:jobName/:id", ms.getAPIJobInvocationOutputStream)
}

// ViewPaths returns the view paths for the management server.
func (ms ManagementServer) ViewPaths() []string {
	return []string{
		"_views/header.html",
		"_views/footer.html",
		"_views/index.html",
		"_views/job.html",
		"_views/invocation.html",
		"_views/partials/job_table.html",
		"_views/partials/job_row.html",
	}
}

// getStatus is mapped to GET /status.json
func (ms ManagementServer) getStatus(r *web.Ctx) web.Result {
	return web.JSON.Result(ms.Cron.Status())
}

// getStatic is mapped to GET /static/*filepath
func (ms ManagementServer) getStatic(r *web.Ctx) web.Result {
	path, err := r.RouteParam("filepath")
	if err != nil {
		web.Text.NotFound()
	}
	path = filepath.Join("_static", path)
	file, err := static.GetBinaryAsset(path)
	if err == os.ErrNotExist {
		return web.Text.NotFound()
	}
	contents, err := file.Contents()
	if err != nil {
		return web.Text.InternalError(err)
	}
	http.ServeContent(r.Response, r.Request, path, time.Unix(file.ModTime, 0), bytes.NewReader(contents))
	return nil
}

//
// api or view routes
//

// getIndex is mapped to GET /
func (ms ManagementServer) getIndex(r *web.Ctx) web.Result {
	r.State.Set("show-job-history-link", true)
	return r.Views.View("index", ms.Cron.Status().Jobs)
}

// getIndex is mapped to GET /search?selector=<SELECTOR>
func (ms ManagementServer) getSearch(r *web.Ctx) web.Result {
	selectorParam := web.StringValue(r.QueryValue("selector"))
	if selectorParam == "" {
		return web.RedirectWithMethod("GET", "/")
	}
	sel, err := selector.Parse(selectorParam)
	if err != nil {
		return r.Views.BadRequest(err)
	}
	r.State.Set("selector", sel.String())

	status := ms.Cron.Status()
	status.Jobs = ms.filterJobSchedulers(status.Jobs, func(js cron.JobSchedulerStatus) bool {
		return sel.Matches(js.Labels)
	})
	return r.Views.View("index", status.Jobs)
}

// getPause is mapped to GET /pause
func (ms ManagementServer) getPause(r *web.Ctx) web.Result {
	if err := ms.Cron.Pause(); err != nil {
		return r.Views.BadRequest(err)
	}
	return web.RedirectWithMethod("GET", "/")
}

// getResume is mapped to GET /resume
func (ms ManagementServer) getResume(r *web.Ctx) web.Result {
	if err := ms.Cron.Resume(); err != nil {
		return r.Views.BadRequest(err)
	}
	return web.RedirectWithMethod("GET", "/")
}

// getJob is mapped to GET /job/:jobName
func (ms ManagementServer) getJob(r *web.Ctx) web.Result {
	jobName, err := r.RouteParam("jobName")
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	job, err := ms.Cron.Job(jobName)
	if err != nil || job == nil {
		return r.Views.NotFound()
	}
	return r.Views.View("job", job)
}

// getJobRun is mapped to GET /job.run/:jobName
func (ms ManagementServer) getJobRun(r *web.Ctx) web.Result {
	jobName, err := r.RouteParam("jobName")
	if err != nil {
		return r.Views.BadRequest(err)
	}
	if err := ms.Cron.RunJob(jobName); err != nil {
		return r.Views.BadRequest(err)
	}
	ji, err := ms.Cron.WaitJobScheduled(r.Context(), jobName)
	if err != nil {
		return r.Views.InternalError(err)
	}
	return web.RedirectWithMethodf("GET", "/job.invocation/%s/%s", jobName, ji.ID)
}

// getJobEnable is mapped to GET /job.enable/:jobName
func (ms ManagementServer) getJobEnable(r *web.Ctx) web.Result {
	jobName, err := r.RouteParam("jobName")
	if err != nil {
		return r.Views.BadRequest(err)
	}
	if err := ms.Cron.EnableJobs(jobName); err != nil {
		return r.Views.BadRequest(err)
	}
	return web.RedirectWithMethod("GET", "/")
}

// getJobDisable is mapped to GET /job.disable/:jobName
func (ms ManagementServer) getJobDisable(r *web.Ctx) web.Result {
	jobName, err := r.RouteParam("jobName")
	if err != nil {
		return r.Views.BadRequest(err)
	}
	if err := ms.Cron.DisableJobs(jobName); err != nil {
		return r.Views.BadRequest(err)
	}
	return web.RedirectWithMethod("GET", "/")
}

// getJobCancel is mapped to GET /job.cancel;/:jobName
func (ms ManagementServer) getJobCancel(r *web.Ctx) web.Result {
	jobName, err := r.RouteParam("jobName")
	if err != nil {
		return r.Views.BadRequest(err)
	}
	err = ms.Cron.CancelJob(jobName)
	if err != nil {
		return r.Views.BadRequest(err)
	}
	return web.RedirectWithMethod("GET", "/")
}

// getJobInvocation is mapped to GET /job.invocation/:jobName/:id
func (ms ManagementServer) getJobInvocation(r *web.Ctx) web.Result {
	invocation, result := ms.getRequestJobInvocation(r, r.Views)
	if result != nil {
		return result
	}
	return r.Views.View("invocation", invocation)
}

// getAPIJobs is mapped to GET /api/jobs
func (ms ManagementServer) getAPIJobs(r *web.Ctx) web.Result {
	return web.JSON.Result(ms.Cron.Status().Jobs)
}

// getAPIJobs is mapped to GET /api/jobs.running
func (ms ManagementServer) getAPIJobsRunning(r *web.Ctx) web.Result {
	return web.JSON.Result(ms.Cron.Status().Running)
}

// postAPIPause is mapped to POST /api/pause
func (ms ManagementServer) postAPIPause(r *web.Ctx) web.Result {
	if err := ms.Cron.Pause(); err != nil {
		return r.Views.BadRequest(err)
	}
	return web.JSON.OK()
}

// postAPIResume is mapped to POST /api/resume
func (ms ManagementServer) postAPIResume(r *web.Ctx) web.Result {
	if err := ms.Cron.Resume(); err != nil {
		return r.Views.BadRequest(err)
	}
	return web.JSON.OK()
}

// getAPIJob is mapped to GET /api/job/:jobName
func (ms ManagementServer) getAPIJob(r *web.Ctx) web.Result {
	jobName, err := r.RouteParam("jobName")
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	job, err := ms.Cron.Job(jobName)
	if err != nil || job == nil {
		return web.JSON.NotFound()
	}
	return web.JSON.Result(job.Status())
}

// getAPIJobsStats is mapped to GET /api/jobs.stats
func (ms ManagementServer) getAPIJobsStats(r *web.Ctx) web.Result {
	output := make(map[string]cron.JobStats)
	for _, job := range ms.Cron.Jobs {
		output[job.Name()] = job.Stats()
	}
	return web.JSON.Result(output)
}

// getAPIJobStats is mapped to GET /api/job.stats/:jobName
func (ms ManagementServer) getAPIJobStats(r *web.Ctx) web.Result {
	jobName, err := r.RouteParam("jobName")
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	job, err := ms.Cron.Job(jobName)
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	if err := ms.Cron.RunJob(jobName); err != nil {
		return web.JSON.BadRequest(err)
	}
	return web.JSON.Result(job.Stats())
}

// postAPIJobRun is mapped to POST /api/job.run/:jobName
func (ms ManagementServer) postAPIJobRun(r *web.Ctx) web.Result {
	jobName, err := r.RouteParam("jobName")
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	if err := ms.Cron.RunJob(jobName); err != nil {
		return web.JSON.BadRequest(err)
	}
	ji, err := ms.Cron.WaitJobScheduled(r.Context(), jobName)
	if err != nil {
		return r.Views.InternalError(err)
	}
	return web.JSON.Result(ji)
}

// postAPIJobCancel is mapped to POST /api/job.cancel/:jobName
func (ms ManagementServer) postAPIJobCancel(r *web.Ctx) web.Result {
	jobName, err := r.RouteParam("jobName")
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	if err := ms.Cron.CancelJob(jobName); err != nil {
		return web.JSON.BadRequest(err)
	}
	return web.JSON.OK()
}

// postAPIJobDisable is mapped to POST /api/job.disable/:jobName
func (ms ManagementServer) postAPIJobDisable(r *web.Ctx) web.Result {
	jobName, err := r.RouteParam("jobName")
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	if err := ms.Cron.DisableJobs(jobName); err != nil {
		return web.JSON.BadRequest(err)
	}
	return web.JSON.OK()
}

// postAPIJobEnable is mapped to POST /api/job.enable/:jobName
func (ms ManagementServer) postAPIJobEnable(r *web.Ctx) web.Result {
	jobName, err := r.RouteParam("jobName")
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	if err := ms.Cron.EnableJobs(jobName); err != nil {
		return web.JSON.BadRequest(err)
	}
	return web.JSON.Result(fmt.Sprintf("%s enabled", jobName))
}

// getAPIJobInvocation is mapped to GET /api/job.invocation/:jobName/:id
func (ms ManagementServer) getAPIJobInvocation(r *web.Ctx) web.Result {
	invocation, result := ms.getRequestJobInvocation(r, web.JSON)
	if result != nil {
		return result
	}
	return web.JSON.Result(invocation)
}

func (ms ManagementServer) getAPIJobInvocationOutput(r *web.Ctx) web.Result {
	invocation, result := ms.getRequestJobInvocation(r, web.JSON)
	if result != nil {
		return result
	}
	lines := append(invocation.Output.Lines)
	if !invocation.Output.Working.Timestamp.IsZero() {
		lines = append(lines, invocation.Output.Working)
	}
	if afterNanos, _ := web.Int64Value(r.QueryValue("afterNanos")); afterNanos > 0 {
		afterTS := time.Unix(0, afterNanos)
		lines = cron.FilterOutputLines(lines, func(l cron.OutputLine) bool {
			return l.Timestamp.After(afterTS)
		})
	}
	return web.JSON.Result(map[string]interface{}{
		"serverTimeNanos": time.Now().UTC().UnixNano(),
		"lines":           lines,
	})
}

func (ms ManagementServer) getAPIJobInvocationOutputStream(r *web.Ctx) web.Result {
	log := r.App.Log

	invocation, result := ms.getRequestJobInvocation(r, web.JSON)
	if result != nil {
		return result
	}
	if !ms.Cron.IsJobInvocationRunning(invocation.JobName, invocation.ID) {
		return nil
	}

	es := webutil.EventSource{Output: r.Response}
	if err := es.StartSession(); err != nil {
		logger.MaybeError(log, err)
		return nil
	}
	listenerID := uuid.V4().String()

	shouldClose := make(chan struct{})
	invocation.OutputListeners.Add(listenerID, func(l cron.OutputLine) {
		if err := es.Data(string(l.Data)); err != nil {
			logger.MaybeError(log, err)
			if shouldClose != nil {
				close(shouldClose)
				shouldClose = nil
			}
		}
	})
	defer func() { invocation.OutputListeners.Remove(listenerID) }()

	updateTick := time.Tick(500 * time.Millisecond)
	for {
		select {
		case <-shouldClose:
			if err := es.EventData("complete", string(invocation.State)); err != nil {
				logger.MaybeError(log, err)
				return nil
			}
			return nil
		case <-updateTick:
			if !ms.Cron.IsJobInvocationRunning(invocation.JobName, invocation.ID) {
				if err := es.EventData("complete", string(invocation.State)); err != nil {
					logger.MaybeError(log, err)
					return nil
				}
				return nil
			}
			if err := es.Ping(); err != nil {
				logger.MaybeError(log, err)
				return nil
			}
			if err := es.EventData("elapsed", fmt.Sprintf("%v", time.Now().UTC().Sub(invocation.Started).Round(time.Millisecond))); err != nil {
				logger.MaybeError(log, err)
				return nil
			}
		}
	}
}

// addContextStateConfig is a middleware that adds the config to a request context's state.
func (ms ManagementServer) addContextStateConfig(action web.Action) web.Action {
	return func(r *web.Ctx) web.Result {
		r.State.Set("config", ms.Config)
		return action(r)
	}
}

// getRequestJobInvocation pulls a job invocation off a request context.
func (ms ManagementServer) getRequestJobInvocation(r *web.Ctx, resultProvider web.ResultProvider) (*cron.JobInvocation, web.Result) {
	jobName, err := r.RouteParam("jobName")
	if err != nil {
		return nil, resultProvider.BadRequest(err)
	}
	job, err := ms.Cron.Job(jobName)
	if err != nil {
		return nil, resultProvider.NotFound()
	}
	invocationID, err := r.RouteParam("id")
	if err != nil {
		return nil, resultProvider.BadRequest(err)
	}

	var invocation *cron.JobInvocation
	if invocationID == "current" && job.Current != nil {
		return job.Current, nil
	}

	invocation = job.GetInvocationByID(invocationID)
	if invocation == nil {
		return nil, resultProvider.NotFound()
	}
	return invocation, nil
}

func (ms ManagementServer) filterJobSchedulers(schedulers []cron.JobSchedulerStatus, predicate func(cron.JobSchedulerStatus) bool) []cron.JobSchedulerStatus {
	var output []cron.JobSchedulerStatus
	for _, js := range schedulers {
		if predicate(js) {
			output = append(output, js)
		}
	}
	return output
}
